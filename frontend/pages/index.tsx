"use client";
import { useEffect, useState } from "react";

export default function Home() {
  const [msg, setMsg] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/api/hello")
      .then((res) => res.json())
      .then((data) => setMsg(data.message));
  }, []);

  return <h1>{msg}</h1>;
}
