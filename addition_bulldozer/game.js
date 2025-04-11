class AdditionBulldozer {
    constructor() {
        this.level = 1;
        this.score = 0;
        this.floors = 0;
        this.correctAnswers = 0;
        this.timeLeft = 5;
        this.maxFloors = 10;
        this.gameOver = false;

        // DOM elements
        this.problemEl = document.getElementById('problem');
        this.answerEl = document.getElementById('answer');
        this.submitBtn = document.getElementById('submit');
        this.feedbackEl = document.getElementById('feedback');
        this.levelEl = document.getElementById('level');
        this.scoreEl = document.getElementById('score');
        this.buildingEl = document.getElementById('building');
        this.floorsEl = document.getElementById('floors');
        this.timerEl = document.getElementById('timer');

        // Event listeners
        this.submitBtn.addEventListener('click', () => this.checkAnswer());
        this.answerEl.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.checkAnswer();
        });

        // Initialize game
        this.generateProblem();
        this.startFloorTimer();
    }

    generateProblem() {
        let num1, num2;
        
        switch(this.level) {
            case 1:
                num1 = Math.floor(Math.random() * 10) + 1;
                num2 = Math.floor(Math.random() * 10) + 1;
                break;
            case 2:
                num1 = Math.floor(Math.random() * 50) + 10;
                num2 = Math.floor(Math.random() * 50) + 10;
                break;
            case 3:
                num1 = Math.floor(Math.random() * 100) + 50;
                num2 = Math.floor(Math.random() * 100) + 50;
                break;
            default:
                num1 = Math.floor(Math.random() * 500) + 100;
                num2 = Math.floor(Math.random() * 500) + 100;
        }

        this.currentProblem = {
            num1,
            num2,
            answer: num1 + num2
        };

        this.problemEl.textContent = `${num1} + ${num2} = ?`;
    }

    checkAnswer() {
        if (this.gameOver) return;

        const userAnswer = parseInt(this.answerEl.value);
        if (isNaN(userAnswer)) {
            this.showFeedback('Please enter a valid number', 'wrong');
            return;
        }

        if (userAnswer === this.currentProblem.answer) {
            this.correctAnswers++;
            this.score += this.level * 10;
            this.removeFloor();
            this.showFeedback('Correct!', 'correct');
            
            if (this.correctAnswers % 4 === 0) {
                this.level++;
                this.levelEl.textContent = this.level;
            }
        } else {
            this.showFeedback('Wrong answer, try again!', 'wrong');
        }

        this.scoreEl.textContent = this.score;
        this.answerEl.value = '';
        this.generateProblem();
    }

    addFloor() {
        if (this.gameOver) return;

        const floor = document.createElement('div');
        floor.className = 'floor new';
        this.buildingEl.appendChild(floor);
        this.floors++;
        this.floorsEl.textContent = this.floors;

        if (this.floors >= this.maxFloors) {
            this.endGame();
        }
    }

    removeFloor() {
        if (this.floors > 0) {
            const floor = this.buildingEl.lastChild;
            floor.classList.add('remove');
            setTimeout(() => {
                this.buildingEl.removeChild(floor);
            }, 500);
            this.floors--;
            this.floorsEl.textContent = this.floors;
        }
    }

    startFloorTimer() {
        const timer = setInterval(() => {
            if (this.gameOver) {
                clearInterval(timer);
                return;
            }

            this.timeLeft--;
            this.timerEl.textContent = this.timeLeft;

            if (this.timeLeft <= 0) {
                this.addFloor();
                this.timeLeft = 5;
            }
        }, 1000);
    }

    showFeedback(message, type) {
        this.feedbackEl.textContent = message;
        this.feedbackEl.className = `feedback ${type}`;
    }

    endGame() {
        this.gameOver = true;
        this.showFeedback(`Game Over! Final Score: ${this.score}`, 'wrong');
        this.submitBtn.disabled = true;
        this.answerEl.disabled = true;
    }
}

// Start the game when the page loads
window.addEventListener('load', () => {
    new AdditionBulldozer();
});
