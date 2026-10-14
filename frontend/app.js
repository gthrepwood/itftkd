const API_BASE = '/api';

// Speech configuration - can be overridden before DOMContentLoaded
const speechConfig = {
    lang: 'hu',
    rate: 0.9,
    pitch: 1,
    volume: 1
};

let quizId = null;
let currentQuestion = null;
let totalQuestions = 0;
let score = 0;
let currentIndex = 0;
const topicLabels = {
    "numbers": "Számok",
    "phrases": "Kifejezések",
    "stances": "Állások",
    "directions": "Irányok",
    "techniques": "Technikák",
    "hand_techniques": "Kéztechnikák",
    "blocks": "Blokkok",
    "kicks": "Rúgások",
    "terms": "Szakkifejezések",
    "principles": "Alapelvek",
    "oath": "Fogadalom",
    "tul": "Formagyakorlatok"
};

async function loadTopics() {
    const res = await fetch(`${API_BASE}/topics`);
    const data = await res.json();
    const select = document.getElementById('topic-select');
    select.innerHTML = '<option value="">Összes téma (vegyes)</option>';
    data.topics.forEach(t => {
        const opt = document.createElement('option');
        opt.value = t;
        opt.textContent = topicLabels[t] || t;
        select.appendChild(opt);
    });
}

function showScreen(screenId) {
    document.getElementById('welcome-screen').classList.add('hidden');
    document.getElementById('quiz-screen').classList.add('hidden');
    document.getElementById('result-screen').classList.add('hidden');
    document.getElementById(screenId).classList.remove('hidden');
}

function startQuiz() {
    const count = parseInt(document.getElementById('question-count').value);
    const topic = document.getElementById('topic-select').value;
    
    fetch(`${API_BASE}/quiz/start`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ questionCount: count, topic: topic })
    })
    .then(res => res.json())
    .then(data => {
        quizId = data.quiz_id;
        totalQuestions = data.total;
        currentIndex = 0;
        score = 0;
        showScreen('quiz-screen');
        updateScore();
        loadQuestion();
    });
}

function loadQuestion() {
    fetch(`${API_BASE}/quiz/${quizId}/question`)
        .then(res => res.json())
        .then(data => {
            if (data.finished) {
                showResults(data.score, data.total);
                return;
            }

            currentIndex = data.index;
            currentQuestion = data.question;
            
            document.getElementById('quiz-category').textContent = topicLabels[currentQuestion.category] || currentQuestion.category;
            document.getElementById('question-text').textContent = currentQuestion.question;
            document.getElementById('progress-bar').style.width = `${(currentIndex / totalQuestions) * 100}%`;
            
            const optionsContainer = document.getElementById('options-container');
            optionsContainer.innerHTML = '';
            document.getElementById('feedback').classList.add('hidden');
            
            currentQuestion.options.forEach((option, index) => {
                const btn = document.createElement('button');
                btn.className = 'option-btn w-full text-left p-4 rounded-xl text-white font-medium flex items-center justify-between';
                btn.innerHTML = `<span class="text-yellow-400 font-bold mr-2">${String.fromCharCode(65 + index)}.</span> <span class="flex-1 text-left">${option}</span> <button class="speak-btn" onclick="event.stopPropagation(); speakText('${option.replace(/'/g, "\\'")}')" aria-label="Felolvasás"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"></polygon><path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"></path></svg></button>`;
                btn.onclick = () => selectAnswer(index, btn);
                optionsContainer.appendChild(btn);
            });
        });
}

function selectAnswer(selectedIndex, buttonElement) {
    const buttons = document.querySelectorAll('.option-btn');
    buttons.forEach(btn => btn.disabled = true);

    document.getElementById('feedback').classList.remove('hidden');

    fetch(`${API_BASE}/quiz/${quizId}/answer`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ answer: selectedIndex })
    })
    .then(res => res.json())
    .then(data => {
        if (data.correct) {
            buttonElement.classList.add('correct');
            score = data.score;
            updateScore();
            document.getElementById('feedback-text').textContent = 'Helyes! Így tovább! 🎉';
            document.getElementById('feedback-text').className = 'text-lg font-medium mb-4 text-green-400';
        } else {
            buttonElement.classList.add('wrong');
            buttons[data.correct_id].classList.add('correct');
            document.getElementById('feedback-text').textContent = 'Helytelen. Ne add fel! 💪';
            document.getElementById('feedback-text').className = 'text-lg font-medium mb-4 text-red-400';
        }
    });
}

function nextQuestion() {
    loadQuestion();
}

function updateScore() {
    document.getElementById('quiz-score').textContent = score;
}

function showResults(finalScore, total) {
    const percentage = (finalScore / total) * 100;
    let emoji, message;

    if (percentage >= 90) {
        emoji = '🏆';
        message = 'Király! Tényleg jól ismered az ITF Taekwondo-t!';
    } else if (percentage >= 70) {
        emoji = '⭐';
        message = 'Nagyon jó! Szilárd az ITF Taekwondo tudásod.';
    } else if (percentage >= 50) {
        emoji = '👍';
        message = 'Szép munka! Gyakorolj, és egyre jobb lesz!';
    } else {
        emoji = '📚';
        message = 'Ne add fel! Minden edzés egy lépés előre.';
    }

    document.getElementById('result-emoji').textContent = emoji;
    document.getElementById('result-message').textContent = message;
    document.getElementById('result-score').textContent = finalScore;
    document.getElementById('result-total').textContent = total;
    
    showScreen('result-screen');
}

function restartQuiz() {
    startQuiz();
}

function goHome() {
    showScreen('welcome-screen');
}

function speakText(text) {
    if ('speechSynthesis' in window) {
        window.speechSynthesis.cancel();
        const utterance = new SpeechSynthesisUtterance(text);
        utterance.lang = speechConfig.lang;
        utterance.rate = speechConfig.rate;
        utterance.pitch = speechConfig.pitch;
        utterance.volume = speechConfig.volume;
        
        const btn = event.target.closest('.speak-btn');
        if (btn) {
            btn.classList.add('speaking');
            utterance.onend = () => btn.classList.remove('speaking');
            utterance.onerror = () => btn.classList.remove('speaking');
        }
        
        window.speechSynthesis.speak(utterance);
    }
}

window.addEventListener('DOMContentLoaded', loadTopics);
