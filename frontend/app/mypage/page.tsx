"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

type Video = {
  id: number;
  title: string;
  thumbnail_url: string;
  created_at: string;
};

export default function MyPage() {
  const router = useRouter();
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(true);
  const [deletingId, setDeletingId] = useState<number | null>(null);

  const fetchMyVideos = () => {
    const token = localStorage.getItem("token");
    if (!token) { router.push("/auth"); return; }

    fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/me/videos`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => (res.ok ? res.json() : []))
      .then((data) => setVideos(Array.isArray(data) ? data : []))
      .catch(() => setVideos([]))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchMyVideos(); }, []);

  const handleDelete = async (videoId: number) => {
    if (!confirm("この動画を削除しますか？")) return;
    setDeletingId(videoId);
    const token = localStorage.getItem("token");
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${videoId}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    });
    if (res.ok) {
      setVideos((prev) => prev.filter((v) => v.id !== videoId));
    } else {
      alert("削除に失敗しました");
    }
    setDeletingId(null);
  };

  return (
    <main style={{ maxWidth: 900, margin: "0 auto", padding: "2rem" }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "2rem" }}>
        <h1 style={{ margin: 0 }}>マイページ</h1>
        <Link href="/videos/new" style={{
          padding: "0.5rem 1rem", background: "#e00", color: "#fff",
          textDecoration: "none", borderRadius: 4,
        }}>
          + 投稿する
        </Link>
      </div>

      {loading ? (
        <p style={{ color: "#888" }}>読み込み中...</p>
      ) : videos.length === 0 ? (
        <div style={{ textAlign: "center", padding: "4rem 0", color: "#888" }}>
          <div style={{ fontSize: "4rem", marginBottom: "1rem" }}>▶</div>
          <p style={{ fontSize: "1.1rem", margin: "0 0 1rem" }}>まだ投稿した動画がありません</p>
          <Link href="/videos/new" style={{ color: "#e00", textDecoration: "none" }}>
            最初の動画を投稿する →
          </Link>
        </div>
      ) : (
        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(260px, 1fr))", gap: "1.5rem" }}>
          {videos.map((video) => (
            <div key={video.id} style={{ border: "1px solid #ddd", borderRadius: 8, overflow: "hidden", background: "#fff" }}>
              <Link href={`/videos/${video.id}`} style={{ textDecoration: "none", color: "inherit" }}>
                <div style={{ background: "#111", height: 150, display: "flex", alignItems: "center", justifyContent: "center" }}>
                  {video.thumbnail_url ? (
                    <img src={video.thumbnail_url} alt={video.title} style={{ width: "100%", height: "100%", objectFit: "cover" }} />
                  ) : (
                    <span style={{ color: "#555", fontSize: "2rem" }}>▶</span>
                  )}
                </div>
                <div style={{ padding: "0.75rem" }}>
                  <p style={{ margin: 0, fontWeight: "bold", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                    {video.title}
                  </p>
                  <p style={{ margin: "0.25rem 0 0", fontSize: "0.8rem", color: "#888" }}>
                    {new Date(video.created_at).toLocaleDateString("ja-JP")}
                  </p>
                </div>
              </Link>
              <div style={{ padding: "0 0.75rem 0.75rem", display: "flex", gap: "0.5rem" }}>
                <Link href={`/videos/${video.id}/edit`} style={{
                  flex: 1, textAlign: "center", padding: "0.4rem", border: "1px solid #ccc",
                  borderRadius: 4, textDecoration: "none", color: "#333", fontSize: "0.85rem",
                }}>
                  編集
                </Link>
                <button
                  onClick={() => handleDelete(video.id)}
                  disabled={deletingId === video.id}
                  style={{
                    flex: 1, padding: "0.4rem", border: "none", borderRadius: 4,
                    background: deletingId === video.id ? "#ccc" : "#e00",
                    color: "#fff", cursor: deletingId === video.id ? "not-allowed" : "pointer",
                    fontSize: "0.85rem",
                  }}
                >
                  {deletingId === video.id ? "削除中..." : "削除"}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </main>
  );
}
