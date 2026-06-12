"use client";

import { useState } from "react";

type Props = {
  videoID: number;
  initialBookmarked: boolean;
};

export default function BookmarkButton({ videoID, initialBookmarked }: Props) {
  const [bookmarked, setBookmarked] = useState(initialBookmarked);

  const handleToggle = async () => {
    const token = localStorage.getItem("token");
    if (!token) return;

    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/api/videos/${videoID}/bookmark`,
      {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      }
    );
    if (!res.ok) return;

    const data = await res.json();
    setBookmarked(data.bookmarked);
  };

  return (
    <button
      onClick={handleToggle}
      aria-label={bookmarked ? "ブックマーク済み" : "ブックマーク"}
      style={{
        background: "none",
        border: "none",
        cursor: "pointer",
        color: bookmarked ? "#e00" : "#aaa",
        padding: "0.4rem",
      }}
    >
      {bookmarked ? (
        <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
          <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z" />
        </svg>
      ) : (
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2z" />
        </svg>
      )}
    </button>
  );
}
