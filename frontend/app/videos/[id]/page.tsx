"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import Link from "next/link";

type Video = {
  id: number;
  user_id: number;
  title: string;
  description: string;
  thumbnail_url: string;
  video_url: string;
  created_at: string;
};

function getMyUserId(): number | null {
  const token = localStorage.getItem("token");
  if (!token) return null;
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.user_id ?? null;
  } catch {
    return null;
  }
}

export default function VideoDetailPage() {
  const { id } = useParams();
  const router = useRouter();
  const [video, setVideo] = useState<Video | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);
  const [myUserId, setMyUserId] = useState<number | null>(null);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    setMyUserId(getMyUserId());
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}`)
      .then((res) => { if (!res.ok) { setNotFound(true); return null; } return res.json(); })
      .then((data) => { if (data) setVideo(data); })
      .catch(() => setNotFound(true))
      .finally(() => setLoading(false));
  }, [id]);

  const handleDelete = async () => {
    if (!confirm("この動画を削除しますか？")) return;
    setDeleting(true);
    const token = localStorage.getItem("token");
    const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}`, {
      method: "DELETE",
      headers: { Authorization: `Bearer ${token}` },
    });
    if (res.ok) {
      router.push("/");
    } else {
      alert("削除に失敗しました");
      setDeleting(false);
    }
  };

  if (loading) return <main style={{ padding: "2rem" }}><p style={{ color: "#888" }}>読み込み中...</p></main>;

  if (notFound || !video) return (
    <main style={{ maxWidth: 800, margin: "0 auto", padding: "2rem" }}>
      <p>動画が見つかりません。</p>
      <Link href="/" style={{ color: "#e00", textDecoration: "none" }}>← 一覧へ戻る</Link>
    </main>
  );

  const isOwner = myUserId !== null && myUserId === video.user_id;

  return (
    <main style={{ maxWidth: 800, margin: "0 auto", padding: "2rem" }}>
      <Link href="/" style={{ color: "#888", textDecoration: "none" }}>← 一覧へ戻る</Link>

      <div style={{ marginTop: "1rem", borderRadius: 8, overflow: "hidden", background: "#111" }}>
        {video.video_url ? (
          <video controls style={{ width: "100%", maxHeight: 450, display: "block" }} poster={video.thumbnail_url || undefined}>
            <source src={video.video_url} type="video/mp4" />
          </video>
        ) : video.thumbnail_url ? (
          <img src={video.thumbnail_url} alt={video.title} style={{ width: "100%", maxHeight: 450, objectFit: "cover", display: "block" }} />
        ) : (
          <div style={{ height: 400, display: "flex", alignItems: "center", justifyContent: "center" }}>
            <span style={{ color: "#555", fontSize: "4rem" }}>▶</span>
          </div>
        )}
      </div>

      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start", marginTop: "1rem" }}>
        <div>
          <h1 style={{ margin: "0 0 0.25rem" }}>{video.title}</h1>
          <p style={{ color: "#888", fontSize: "0.9rem", margin: 0 }}>
            {new Date(video.created_at).toLocaleDateString("ja-JP")}
          </p>
        </div>
        {isOwner && (
          <div style={{ display: "flex", gap: "0.5rem", flexShrink: 0 }}>
            <Link href={`/videos/${video.id}/edit`} style={{
              padding: "0.4rem 1rem", border: "1px solid #ccc", borderRadius: 4,
              textDecoration: "none", color: "#333", fontSize: "0.9rem",
            }}>
              編集
            </Link>
            <button onClick={handleDelete} disabled={deleting} style={{
              padding: "0.4rem 1rem", background: deleting ? "#ccc" : "#e00",
              color: "#fff", border: "none", borderRadius: 4,
              cursor: deleting ? "not-allowed" : "pointer", fontSize: "0.9rem",
            }}>
              {deleting ? "削除中..." : "削除"}
            </button>
          </div>
        )}
      </div>

      {video.description && (
        <p style={{ marginTop: "1rem", lineHeight: 1.7, borderTop: "1px solid #eee", paddingTop: "1rem" }}>
          {video.description}
        </p>
      )}
    </main>
  );
}
