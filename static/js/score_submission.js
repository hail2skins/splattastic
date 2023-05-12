let scoresData = {};
let forceRefresh = false;

function loadScores(diveId) {
    console.log("Loading scores for dive:", diveId);
    const userId = document.querySelector('#scoreForm input[name="userId"]').value;
    const eventId = document.querySelector('#scoreForm input[name="eventId"]').value;

    // If forceRefresh is true, clear the scores data for the dive
    if (forceRefresh) {
        scoresData[diveId] = null;
        forceRefresh = false; // Reset the flag
    }

    if (!scoresData[diveId]) {
        scoresData[diveId] = { scores: [] };
        fetch(`/user/${userId}/event/${eventId}/dive/${diveId}/scores`)
            .then(function (response) {
                if (response.ok) {
                    return response.json();
                } else {
                    throw new Error('Failed to fetch scores.');
                }
            })
            .then(function (data) {
                console.log("Fetched scores:", data);
                // Clear the scores array before adding new scores
                scoresData[diveId].scores = [];
                scoresData[diveId].scores = data.scores;
                updateScores(diveId);

                // Fetch and display the total score if there are 3, 5, or 7 scores
                if ([3, 5, 7].includes(data.scores.length)) {
                    fetch(`/user/${userId}/event/${eventId}/dive/${diveId}/total`, {
                        method: "GET"
                    })
                    .then(response => response.json())
                    .then(data => {
                        // Display the total score
                        const totalScoreElement = document.querySelector(`.dive-${diveId} .total-score`);
                        if (totalScoreElement) {
                            totalScoreElement.textContent = data.score.toFixed(2);
                        } else {
                            const newTotalScoreElement = document.createElement("div");
                            newTotalScoreElement.classList.add("total-score");
                            newTotalScoreElement.textContent = data.score.toFixed(2);
                            document.querySelector(`[data-dive-id="${diveId}"]`).appendChild(newTotalScoreElement);
                        }

                        // Fetch and display the meet score
                        fetch(`/user/${userId}/event/${eventId}/meet_score`, {
                            method: "GET"
                        })
                        .then(response => response.json())
                        .then(data => {
                            // Display the meet score
                            const meetScoreElement = document.querySelector("#meetScore");
                            if (meetScoreElement) {
                                if (data.score === null || data.score === undefined) {
                                    meetScoreElement.textContent = "No dives scored yet";
                                } else {
                                    meetScoreElement.textContent = data.score.toFixed(2);
                                }
                            }
                        });
                    });
                }
            })
            .catch(function (error) {
                console.error('Error fetching scores:', error);
            });
    } else {
        updateScores(diveId);
    }
}

function updateScores(diveId) {
    // First, remove any existing scores from the DOM
    const scoreElements = document.querySelectorAll(`[data-dive-id="${diveId}"] .score`);
    scoreElements.forEach(function (scoreElement) {
        scoreElement.remove();
    });

    // Then add the new scores
    scoresData[diveId].scores.forEach(function (scoreObject, index) {
        const score = scoreObject.score;
        const judge = scoreObject.judge; // Assuming scoreObject contains a judge field

        // Create a new score element and add it to the DOM
        const parentContainer = document.querySelector(`[data-dive-id="${diveId}"][data-score-index="${index + 1}"]`);
        const newScoreElement = document.createElement("div");
        newScoreElement.classList.add("score");  // Add the 'score' class
        newScoreElement.textContent = score.toFixed(2);
        parentContainer.appendChild(newScoreElement);

        // Update and disable the corresponding input field in the modal
        const scoreInput = document.querySelector(`#scoreForm input[name="score${judge}"]`);
        if (scoreInput) {
            scoreInput.value = score.toFixed(2);
        }
    });
}



document.addEventListener("DOMContentLoaded", function() {
    const scoreForm = document.getElementById("scoreForm");
    const scoreModal = document.getElementById("scoreModal");
    const errorAlert = document.getElementById("errorAlert");

    // Handle form submission
    scoreForm.addEventListener("submit", function(event) {
        event.preventDefault();
        console.log("Form submitted"); // Add this line

        // Get the form data
        const formData = new FormData(scoreForm);
        const diveId = parseInt(formData.get("diveId"), 10);
        const userId = parseInt(formData.get("userId"), 10);
        const eventId = parseInt(formData.get("eventId"), 10);

        // Get scores from the form data
        const scores = [];
        for (let i = 1; i <= 7; i++) {
            const score = formData.get("score" + i);
            if (score) {
                // Validation: scores must be between 0 and 10 and in increments of 0.5
                const floatScore = parseFloat(score);
                if (floatScore < 0 || floatScore > 10 || floatScore * 2 % 1 !== 0) {
                    // Invalid score - show error message and return to prevent form submission
                    errorAlert.textContent = "Scores must be between 0 and 10 and in increments of 0.5.";
                    errorAlert.style.display = "block";
                    return;
                }
                scores.push(floatScore);
            }
        }

        // Send a request to the server to save the scores
        
        const url = `/user/${userId}/event/${eventId}/scores`; // Creates or updates through ScoreUpsert in score model
        fetch(url, {
            method: "POST",
            body: JSON.stringify({
                userId: userId,
                eventId: eventId,
                diveId: diveId,
                scores: scores
            }),
            headers: {
                "Content-Type": "application/json"
            }
        }).then(function(response) {
            if (response.ok) {
                return response.json(); // Return the response as a Promise
            } else {
                throw new Error("An error occurred while submitting the score.");
            }
        }).then(function(data) {
            // Close the modal and clear the form
            const modalInstance = bootstrap.Modal.getInstance(scoreModal);
            modalInstance.hide();
            scoreForm.reset();

            // Set the force refresh flag and reload the scores
            forceRefresh = true;
            loadScores(diveId);

            // Update total score in the DOM
            const totalScoreElement = document.querySelector(`.dive-${diveId} .total-score`);
            if (data.total_score) {
                totalScoreElement.textContent = data.total_score.toFixed(2);
            } else {
                totalScoreElement.textContent = ''; // Clear the total score if not available
            }

            scores.forEach(function(score, index) {
                // Select the parent container for the current score index
                const parentContainer = document.querySelector(`[data-dive-id="${diveId}"][data-score-index="${index + 1}"]`);
                
                const scoreElement = parentContainer.querySelector(".score");
                if (scoreElement) {
                    scoreElement.textContent = score.toFixed(2);
                } else {
                    const newScoreElement = document.createElement("div");
                    newScoreElement.classList.add("score");  // Add the 'score' class
                    newScoreElement.textContent = score.toFixed(2);
                    parentContainer.appendChild(newScoreElement);
                }
            });

            // If there are 3, 5, or 7 scores, calculate the total score
            if ([3, 5, 7].includes(scores.length)) {
                fetch(`/user/${userId}/event/${eventId}/dive/${diveId}/total`, {
                    method: "GET"
                })
                .then(response => response.json())
                .then(data => {
                    // Display the total score
                    const totalScoreElement = document.querySelector(`.dive-${diveId} .total-score`);
                    if (totalScoreElement) {
                        totalScoreElement.textContent = data.score.toFixed(2);
                    } else {
                        const newTotalScoreElement = document.createElement("div");
                        newTotalScoreElement.classList.add("total-score");
                        newTotalScoreElement.textContent = data.score.toFixed(2);
                        document.querySelector(`[data-dive-id="${diveId}"]`).appendChild(newTotalScoreElement);
                    }
                });
            }
            

            // Hide the error alert
            errorAlert.style.display = "none";
        }).catch(function(error) {
            // Handle error
            // Display the error message
            errorAlert.textContent = error.message;
            errorAlert.style.display = "block";
        });
    });

// Set the dive ID in the form when the Add Score button is clicked
const addScoreButtons = document.querySelectorAll(".add-score");
addScoreButtons.forEach(function(button) {
    button.addEventListener("click", function() {
        // Clear the form and enable all input fields
        scoreForm.reset();
        for (let i = 1; i <= 7; i++) {
            const scoreInput = document.querySelector(`#scoreForm input[name="score${i}"]`);
        }

        const diveId = button.getAttribute("data-dive-id");
        console.log("Add Score button clicked, diveId:", diveId); // Add this line
        scoreForm.querySelector("#diveId").value = diveId;
        loadScores(diveId); // Refresh the scores for the selected dive
    });
});


    // Load scores for all dives
    document.querySelectorAll('[data-dive-id]').forEach(function (element) {
        const diveId = element.getAttribute('data-dive-id');
        loadScores(diveId);
    });

});
