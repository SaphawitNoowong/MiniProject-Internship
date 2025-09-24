import { useEffect, useState } from "react";

export default function Home() {
  const [users, setUsers] = useState<any[]>([]);

  useEffect(() => {
    fetch("http://localhost:5000/api/users")
      .then((res) => res.json())
      .then((data) => setUsers(data))
      .catch((err) => console.error(err));
  }, []);

  return (
    <div>
      <h1>Users from MongoDB</h1>
      <ul>
        {users.map((u, i) => (
          <li key={i}>{JSON.stringify(u)}</li>
        ))}
      </ul>
    </div>
  );
}
