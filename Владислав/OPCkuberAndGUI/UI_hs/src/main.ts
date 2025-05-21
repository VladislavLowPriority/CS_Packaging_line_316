import express from 'express';
import bodyParser from 'body-parser';
import path from 'path';
import { connectToOpcuaServer, sendCommandToOpcua } from './opcuaClient';
import { WebSocketServer } from 'ws';
import { readOpcuaStates } from './opcuaClient';
const app = express();
const PORT = 3000;
const publicPath = path.join(__dirname, '../public');
app.use(express.static(publicPath));
console.log('Serving static files from:', publicPath);
// Middleware
app.use(express.json());

app.get('/', (req, res) => {
    res.sendFile(path.join(publicPath, 'index.html'));
});


//webSocket сервер:
const wss = new WebSocketServer({ port: 8080 });

setInterval(async () => {
    try {
        const state = await readOpcuaStates();
        const data = JSON.stringify(state);
        wss.clients.forEach(client => {
            if (client.readyState === 1) {
                client.send(data);
            }
        });
    } catch (err) {
        console.error("OPC UA read error:", err);
    }
}, 200);

// POST /opcua-command
app.post('/opcua-command', async (req, res) => {
    const { tag, value } = req.body;

    try {
        await sendCommandToOpcua(tag, value);
        res.json({ success: true });
    } catch (err) {
        console.error("Error sending OPC UA command:", err);
        res.status(500).json({ success: false, error: err});
    }
});

// Start
app.listen(PORT, async () => {
    console.log(`Server running at http://localhost:${PORT}`);
    try {
        await connectToOpcuaServer();
    } catch (err) {
        console.error('Failed to connect to OPC UA server:', err);
        process.exit(1);
    }
});
