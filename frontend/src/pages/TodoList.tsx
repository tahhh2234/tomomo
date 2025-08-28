/* eslint-disable @typescript-eslint/no-unused-vars */
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
  const [error, setError] = useState("");

  const fetchTasks = async () => {
    try {
      const res = await getTasks();
      setTasks(res.data);
    } catch (err) {
      setError("Can not load the tasks.");
    }
  };

  useEffect(() => {
    fetchTasks();
  }, []);

  const handleCreate = async () => {
    if (!newTitle) return;
    try {
      await createTask({ title: newTitle });
      setNewTitle("");
      fetchTasks();
    } catch (err) {
      setError("Can not create the task.");
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteTask(id);
      fetchTasks();
    } catch (err) {
      setError("Can not delete the task.");
    }
  };

  return (
    <div style={{ padding: "20px" }}>
      <h2>Todo List</h2>
      {error && <p style={{ color: "red" }}>{error}</p>}

      <div>
        <input
          value={newTitle}
          onChange={e => setNewTitle(e.target.value)}
          placeholder="New task"
        />
        <button onClick={handleCreate}>Add</button>
      </div>

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
