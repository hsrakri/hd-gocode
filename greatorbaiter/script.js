// Game state
let gameState = {
    playerName: '',
    score: 0,
    level: 1,
    timer: 5,
    correctAnswers: 0,
    timerInterval: null,
    leaderboard: []
};

// DOM elements
const welcomeScreen = document.getElementById('welcome-screen');
const gameScreen = document.getElementById('game-screen');
const gameOverScreen = document.getElementById('game-over-screen');
const playerNameInput = document.getElementById('player-name');
const startGameButton = document.getElementById('start-game');
const playAgainButton = document.getElementById('play-again');
const scoreElement = document.getElementById('score');
const timerElement = document.getElementById('timer');
const levelElement = document.getElementById('level');
const finalScoreElement = document.getElementById('final-score');
const leftNumber = document.getElementById('left-number');
const rightNumber = document.getElementById('right-number');
const leaderboardList = document.getElementById('leaderboard-list');

// Event listeners
startGameButton.addEventListener('click', startGame);
playAgainButton.addEventListener('click', resetGame);
leftNumber.addEventListener('click', () => checkAnswer('left'));
rightNumber.addEventListener('click', () => checkAnswer('right'));

// Initialize leaderboard from localStorage
function initLeaderboard() {
    const savedLeaderboard = localStorage.getItem('greatorbaiterLeaderboard');
    if (savedLeaderboard) {
        gameState.leaderboard = JSON.parse(savedLeaderboard);
    }
}

// Start game
function startGame() {
    if (!playerNameInput.value.trim()) {
        alert('Please enter your name!');
        return;
    }
    
    gameState.playerName = playerNameInput.value.trim();
    welcomeScreen.classList.add('hidden');
    gameScreen.classList.remove('hidden');
    
    generateNewNumbers();
    startTimer();
}

// Generate new numbers based on level
function generateNewNumbers() {
    const maxNumber = Math.pow(10, gameState.level);
    const num1 = Math.floor(Math.random() * maxNumber);
    let num2;
    do {
        num2 = Math.floor(Math.random() * maxNumber);
    } while (num2 === num1);

    // Randomly decide which number goes to which side
    if (Math.random() < 0.5) {
        leftNumber.querySelector('.number').textContent = num1;
        rightNumber.querySelector('.number').textContent = num2;
    } else {
        leftNumber.querySelector('.number').textContent = num2;
        rightNumber.querySelector('.number').textContent = num1;
    }
}

// Timer function
function startTimer() {
    gameState.timer = 5;
    timerElement.textContent = gameState.timer;
    
    if (gameState.timerInterval) {
        clearInterval(gameState.timerInterval);
    }
    
    gameState.timerInterval = setInterval(() => {
        gameState.timer--;
        timerElement.textContent = gameState.timer;
        
        if (gameState.timer <= 0) {
            endGame();
        }
    }, 1000);
}

// Check answer
function checkAnswer(side) {
    const leftValue = parseInt(leftNumber.querySelector('.number').textContent);
    const rightValue = parseInt(rightNumber.querySelector('.number').textContent);
    const isCorrect = (side === 'left' && leftValue > rightValue) || 
                     (side === 'right' && rightValue > leftValue);

    if (isCorrect) {
        gameState.score += gameState.level * 10;
        gameState.correctAnswers++;
        scoreElement.textContent = gameState.score;
        
        // Level up after every 2 correct answers
        if (gameState.correctAnswers % 2 === 0) {
            gameState.level++;
            levelElement.textContent = gameState.level;
        }
        
        generateNewNumbers();
        startTimer();
    } else {
        endGame();
    }
}

// End game
function endGame() {
    clearInterval(gameState.timerInterval);
    gameScreen.classList.add('hidden');
    gameOverScreen.classList.remove('hidden');
    finalScoreElement.textContent = gameState.score;
    
    // Update leaderboard
    gameState.leaderboard.push({
        name: gameState.playerName,
        score: gameState.score
    });
    
    // Sort leaderboard by score (descending)
    gameState.leaderboard.sort((a, b) => b.score - a.score);
    
    // Keep only top 10 scores
    gameState.leaderboard = gameState.leaderboard.slice(0, 10);
    
    // Save to localStorage
    localStorage.setItem('greatorbaiterLeaderboard', JSON.stringify(gameState.leaderboard));
    
    // Display leaderboard
    displayLeaderboard();
}

// Display leaderboard
function displayLeaderboard() {
    leaderboardList.innerHTML = '';
    gameState.leaderboard.forEach((entry, index) => {
        const item = document.createElement('div');
        item.className = 'leaderboard-item';
        item.innerHTML = `
            <span>${index + 1}. ${entry.name}</span>
            <span>${entry.score}</span>
        `;
        leaderboardList.appendChild(item);
    });
}

// Reset game
function resetGame() {
    gameState.score = 0;
    gameState.level = 1;
    gameState.correctAnswers = 0;
    scoreElement.textContent = '0';
    levelElement.textContent = '1';
    
    gameOverScreen.classList.add('hidden');
    welcomeScreen.classList.remove('hidden');
    playerNameInput.value = '';
}

// Initialize the game
initLeaderboard(); 