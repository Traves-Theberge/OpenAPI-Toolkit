import express from 'express';
import { Request, Response } from 'express';

const app = express();
const PORT = 3000;

app.use(express.json());

interface User {
  id: number;
  name: string;
  email: string;
}

let users: User[] = [];
let nextId = 1;

app.get('/health', (req: Request, res: Response) => {
  res.json({ status: 'ok' });
});

app.post('/users', (req: Request, res: Response) => {
  const { name, email } = req.body;

  if (!name || !email) {
    return res.status(400).json({ error: 'Name and email are required' });
  }

  const newUser: User = { id: nextId++, name, email };
  users.push(newUser);

  res.status(201).json(newUser);
});

app.get('/users', (req: Request, res: Response) => {
  res.json(users);
});

app.listen(PORT, () => {
  console.log(`Sample API server running on http://localhost:${PORT}`);
});