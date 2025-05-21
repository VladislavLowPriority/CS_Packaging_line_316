// Current state
let states = {
    start: false,
    stop: false,
    reset: false
};
//webSocket
const ws = new WebSocket('ws://localhost:8080');
ws.onmessage = (event) => {
    const newStates = JSON.parse(event.data);
    states = { ...states, ...newStates };
    updateButtons();
};

// DOM elements
const startBtn = document.getElementById('startBtn');
const stopBtn = document.getElementById('stopBtn');
const resetBtn = document.getElementById('resetBtn');
const statusMessage = document.getElementById('statusMessage');

// Update button states and appearance
function updateButtons() {
    // Start button
    startBtn.classList.toggle('active', states.start);
    
    // Stop button
    stopBtn.classList.toggle('active', states.stop);
    resetBtn.classList.toggle('activeToStart', states.stop);
    
    // Reset button
    resetBtn.classList.toggle('active', states.reset);
    
    // Управление disabled состояниями
    resetBtn.disabled = !states.stop;
    
    updateStatusMessage();
}

// Update status message
function updateStatusMessage() {
    let message = "Status: ";
    if (states.start) {
        message += "System Started";
    } else if (states.stop) {
        if (states.reset) {
            message += "System Stopped and Returning to Start";
        } else {
            message += "System Stopped";
        }
    } else if (states.reset) {
        message += "Returning to Start";
    } else {
        message += "Ready";
    }
    statusMessage.textContent = message;
}

// Function to send command to backend
async function sendCommand(tag, value) {
    try {
        const response = await fetch('/opcua-command', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ tag, value }),
        });

        const data = await response.json();
        
        if (!data.success) {
            console.error('Error from server:', data.error);
            // Revert state if command failed
            states[tag.split('_')[0]] = !value;
            updateButtons();
        }
    } catch (error) {
        console.error('Network error:', error);
        // Revert state if network error
        states[tag.split('_')[0]] = !value;
        updateButtons();
    }
}

// Button click handlers
startBtn.addEventListener('click', () => {
    const newState = !states.start;
    states.start = newState;
    if (states.start) {
        states.stop = false;
        states.reset = false;
    }
    updateButtons();
    sendCommand('start_hs', newState);
    sendCommand('stop_hs', false);
    sendCommand('back_to_start', false);
});

stopBtn.addEventListener('click', () => {
    const newState = !states.stop;
    states.stop = newState;
    if (states.stop) {
        states.start = false;
        states.reset = false;
    }
    updateButtons();
    sendCommand('stop_hs', newState);
    sendCommand('start_hs', false);
    sendCommand('back_to_start', false);
});

resetBtn.addEventListener('click', () => {
    if (states.stop) {
        const newState = !states.reset;
        states.reset = newState;
        if (states.reset) {
            states.start = false;
            states.stop = false;
        }
        updateButtons();
        sendCommand('start_hs', false);
        sendCommand('stop_hs', false);
        sendCommand('back_to_start', newState);  
    }
});

// Initialize buttons
updateButtons();