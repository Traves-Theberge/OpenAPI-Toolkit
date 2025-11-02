"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const app = (0, express_1.default)();
const PORT = 3000;
app.use(express_1.default.json());
let users = [];
let nextId = 1;
app.get('/health', (req, res) => {
    res.json({ status: 'ok' });
});
app.post('/users', (req, res) => {
    const { name, email } = req.body;
    if (!name || !email) {
        return res.status(400).json({ error: 'Name and email are required' });
    }
    const newUser = { id: nextId++, name, email };
    users.push(newUser);
    res.status(201).json(newUser);
});
app.get('/users', (req, res) => {
    res.json(users);
});
app.listen(PORT, () => {
    console.log(`Sample API server running on http://localhost:${PORT}`);
});
//# sourceMappingURL=server.js.map