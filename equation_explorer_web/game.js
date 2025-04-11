class MazeCell {
    constructor(x, y) {
        this.x = x;
        this.y = y;
        this.walls = { top: true, right: true, bottom: true, left: true };
        this.visited = false;
        this.hasEquation = false;
        this.equation = null;
    }
}

class Game {
    constructor() {
        this.canvas = document.getElementById('gameCanvas');
        this.ctx = this.canvas.getContext('2d');
        this.score = 0;
        this.level = 1;
        this.timeLeft = 300; // 5 minutes
        this.gameStarted = false;
        this.maze = [];
        this.cellSize = 40;
        this.player = { x: 0, y: 0 };
        this.equations = [];
        this.currentEquation = null;
        
        // Resize canvas
        this.resizeCanvas();
        window.addEventListener('resize', () => this.resizeCanvas());
        
        // Event listeners
        document.getElementById('startButton').addEventListener('click', () => this.startGame());
        document.getElementById('restartButton').addEventListener('click', () => this.restartGame());
        document.getElementById('submit').addEventListener('click', () => this.checkAnswer());
        document.addEventListener('keydown', (e) => this.handleKeyPress(e));
        
        // Initialize start screen
        document.getElementById('startScreen').classList.remove('hidden');
        document.getElementById('gameOverScreen').classList.add('hidden');
    }

    resizeCanvas() {
        this.canvas.width = this.canvas.offsetWidth;
        this.canvas.height = this.canvas.offsetHeight;
        if (this.gameStarted) {
            this.draw();
        }
    }

    startGame() {
        this.gameStarted = true;
        document.getElementById('startScreen').classList.add('hidden');
        this.initializeLevel();
        this.gameLoop();
        this.startTimer();
    }

    restartGame() {
        this.score = 0;
        this.level = 1;
        this.timeLeft = 300;
        document.getElementById('gameOverScreen').classList.add('hidden');
        this.startGame();
    }

    initializeLevel() {
        const mazeSize = 8 + Math.floor(this.level / 2); // Maze gets bigger with each level
        this.maze = [];
        for (let y = 0; y < mazeSize; y++) {
            this.maze[y] = [];
            for (let x = 0; x < mazeSize; x++) {
                this.maze[y][x] = new MazeCell(x, y);
            }
        }
        
        this.generateMaze(0, 0);
        this.player = { x: 0, y: 0 };
        this.placeEquations();
        this.updateStats();
    }

    generateMaze(x, y) {
        const cell = this.maze[y][x];
        cell.visited = true;
        
        const directions = this.shuffleArray([
            { dx: 0, dy: -1, wall: 'top' },
            { dx: 1, dy: 0, wall: 'right' },
            { dx: 0, dy: 1, wall: 'bottom' },
            { dx: -1, dy: 0, wall: 'left' }
        ]);
        
        for (const dir of directions) {
            const newX = x + dir.dx;
            const newY = y + dir.dy;
            
            if (this.isValidCell(newX, newY) && !this.maze[newY][newX].visited) {
                cell.walls[dir.wall] = false;
                const oppositeWall = this.getOppositeWall(dir.wall);
                this.maze[newY][newX].walls[oppositeWall] = false;
                this.generateMaze(newX, newY);
            }
        }
    }

    isValidCell(x, y) {
        return x >= 0 && x < this.maze[0].length && y >= 0 && y < this.maze.length;
    }

    getOppositeWall(wall) {
        const opposites = { top: 'bottom', right: 'left', bottom: 'top', left: 'right' };
        return opposites[wall];
    }

    shuffleArray(array) {
        for (let i = array.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [array[i], array[j]] = [array[j], array[i]];
        }
        return array;
    }

    placeEquations() {
        const numEquations = 3 + this.level; // More equations with each level
        const cells = [];
        
        // Get all available cells except start and end
        for (let y = 0; y < this.maze.length; y++) {
            for (let x = 0; x < this.maze[0].length; x++) {
                if ((x !== 0 || y !== 0) && (x !== this.maze[0].length - 1 || y !== this.maze.length - 1)) {
                    cells.push({ x, y });
                }
            }
        }
        
        this.shuffleArray(cells);
        for (let i = 0; i < numEquations && i < cells.length; i++) {
            const cell = this.maze[cells[i].y][cells[i].x];
            cell.hasEquation = true;
            cell.equation = this.generateEquation();
        }
    }

    generateEquation() {
        const operations = ['+', '-', '*'];
        const operation = operations[Math.floor(Math.random() * operations.length)];
        let a, b, answer;
        
        switch (operation) {
            case '+':
                a = Math.floor(Math.random() * 20) + 1;
                b = Math.floor(Math.random() * 20) + 1;
                answer = a + b;
                return { text: `x + ${b} = ${answer}`, answer: a };
            case '-':
                answer = Math.floor(Math.random() * 20) + 1;
                b = Math.floor(Math.random() * 20) + 1;
                a = answer + b;
                return { text: `x - ${b} = ${answer}`, answer: a };
            case '*':
                a = Math.floor(Math.random() * 12) + 1;
                b = Math.floor(Math.random() * 12) + 1;
                answer = a * b;
                return { text: `${a}x = ${answer}`, answer: b };
        }
    }

    handleKeyPress(e) {
        if (!this.gameStarted || this.currentEquation) return;
        
        const key = e.key;
        const currentCell = this.maze[this.player.y][this.player.x];
        let newX = this.player.x;
        let newY = this.player.y;
        
        switch (key) {
            case 'ArrowUp':
                if (!currentCell.walls.top) newY--;
                break;
            case 'ArrowRight':
                if (!currentCell.walls.right) newX++;
                break;
            case 'ArrowDown':
                if (!currentCell.walls.bottom) newY++;
                break;
            case 'ArrowLeft':
                if (!currentCell.walls.left) newX--;
                break;
        }
        
        if (this.isValidCell(newX, newY)) {
            this.player.x = newX;
            this.player.y = newY;
            
            const newCell = this.maze[newY][newX];
            if (newCell.hasEquation) {
                this.currentEquation = newCell.equation;
                document.getElementById('equation').textContent = `Solve: ${this.currentEquation.text}`;
                document.getElementById('answer').value = '';
                document.getElementById('answer').focus();
            }
            
            if (newX === this.maze[0].length - 1 && newY === this.maze.length - 1) {
                this.completeLevel();
            }
            
            this.draw();
        }
    }

    checkAnswer() {
        if (!this.currentEquation) return;
        
        const userAnswer = parseInt(document.getElementById('answer').value);
        const correctAnswer = this.currentEquation.answer;
        
        if (userAnswer === correctAnswer) {
            this.score += 100 * this.level;
            document.getElementById('message').textContent = 'Correct!';
            this.maze[this.player.y][this.player.x].hasEquation = false;
            this.currentEquation = null;
            this.updateStats();
        } else {
            document.getElementById('message').textContent = 'Try again!';
            this.score -= 10;
            this.updateStats();
        }
    }

    completeLevel() {
        this.level++;
        this.score += 500 * this.level;
        this.timeLeft += 60; // Bonus time for completing level
        this.updateStats();
        this.initializeLevel();
    }

    startTimer() {
        this.timer = setInterval(() => {
            this.timeLeft--;
            this.updateStats();
            
            if (this.timeLeft <= 0) {
                this.gameOver();
            }
        }, 1000);
    }

    gameOver() {
        clearInterval(this.timer);
        this.gameStarted = false;
        document.getElementById('finalScore').textContent = `Final Score: ${this.score}`;
        document.getElementById('gameOverScreen').classList.remove('hidden');
    }

    updateStats() {
        document.getElementById('score').textContent = `Score: ${this.score}`;
        document.getElementById('level').textContent = `Level: ${this.level}`;
        document.getElementById('timer').textContent = `Time: ${Math.floor(this.timeLeft / 60)}:${(this.timeLeft % 60).toString().padStart(2, '0')}`;
    }

    draw() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        const cellSize = Math.min(
            (this.canvas.width - 20) / this.maze[0].length,
            (this.canvas.height - 20) / this.maze.length
        );
        
        const offsetX = (this.canvas.width - cellSize * this.maze[0].length) / 2;
        const offsetY = (this.canvas.height - cellSize * this.maze.length) / 2;
        
        // Draw maze
        this.ctx.strokeStyle = '#4CAF50';
        this.ctx.lineWidth = 2;
        
        for (let y = 0; y < this.maze.length; y++) {
            for (let x = 0; x < this.maze[0].length; x++) {
                const cell = this.maze[y][x];
                const cellX = offsetX + x * cellSize;
                const cellY = offsetY + y * cellSize;
                
                if (cell.walls.top) {
                    this.ctx.beginPath();
                    this.ctx.moveTo(cellX, cellY);
                    this.ctx.lineTo(cellX + cellSize, cellY);
                    this.ctx.stroke();
                }
                if (cell.walls.right) {
                    this.ctx.beginPath();
                    this.ctx.moveTo(cellX + cellSize, cellY);
                    this.ctx.lineTo(cellX + cellSize, cellY + cellSize);
                    this.ctx.stroke();
                }
                if (cell.walls.bottom) {
                    this.ctx.beginPath();
                    this.ctx.moveTo(cellX, cellY + cellSize);
                    this.ctx.lineTo(cellX + cellSize, cellY + cellSize);
                    this.ctx.stroke();
                }
                if (cell.walls.left) {
                    this.ctx.beginPath();
                    this.ctx.moveTo(cellX, cellY);
                    this.ctx.lineTo(cellX, cellY + cellSize);
                    this.ctx.stroke();
                }
                
                // Draw equation marker
                if (cell.hasEquation) {
                    this.ctx.fillStyle = '#FFD700';
                    this.ctx.beginPath();
                    this.ctx.arc(
                        cellX + cellSize / 2,
                        cellY + cellSize / 2,
                        cellSize / 4,
                        0,
                        Math.PI * 2
                    );
                    this.ctx.fill();
                }
            }
        }
        
        // Draw player
        this.ctx.fillStyle = '#4CAF50';
        this.ctx.beginPath();
        this.ctx.arc(
            offsetX + this.player.x * cellSize + cellSize / 2,
            offsetY + this.player.y * cellSize + cellSize / 2,
            cellSize / 3,
            0,
            Math.PI * 2
        );
        this.ctx.fill();
        
        // Draw exit
        this.ctx.fillStyle = '#FF4444';
        this.ctx.beginPath();
        this.ctx.arc(
            offsetX + (this.maze[0].length - 1) * cellSize + cellSize / 2,
            offsetY + (this.maze.length - 1) * cellSize + cellSize / 2,
            cellSize / 4,
            0,
            Math.PI * 2
        );
        this.ctx.fill();
    }
}

// Start the game when the page loads
window.addEventListener('load', () => {
    new Game();
}); 