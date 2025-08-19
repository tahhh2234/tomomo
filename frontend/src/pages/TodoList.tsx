import { useEffect, useState } from "react";
import { getTasks, createTask, deleteTask } from "../api/api";

type Task = {
  id: number;
  title: string;
  priority: number;
  user_id: number;
};

export default function TodoList() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [newTitle, setNewTitle] = useState("");

  const fetchTasks = async () => {
    const res = await getTasks();
    setTasks(res.data);
  };

  useEffect(() => {
    fetchTasks();
  }, []);

  const handleCreate = async () => {
    if (!newTitle) return;
    await createTask({ title: newTitle });
    setNewTitle("");
    fetchTasks();
  };

  const handleDelete = async (id: number) => {
    await deleteTask(id);
    fetchTasks();
  };

  return (
    <div>
      <h2>Todo List</h2>
      <input value={newTitle} onChange={e => setNewTitle(e.target.value)} placeholder="New task" />
      <button onClick={handleCreate}>Add</button>
      <ul>
        {tasks.map(task => (
          <li key={task.id}>
            {task.title} (priority: {task.priority})
            <button onClick={() => handleDelete(task.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}
