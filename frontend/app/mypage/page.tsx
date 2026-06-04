"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";

type Video = {
  id: number;
  title: string;
  thumbnail_url: string;
  status: string;
  created_at: string;
};

const statusLabel: Record<string, { text: string; color: string; bg: string }> = {
  pending:    { text: "変換待ち",  color: "#b45309", bg: "#fef3c7" },
  processing: { text: "変換中",    color: "#1d4ed8", bg: "#dbeafe" },
  done:       { text: "公開中",    color: "#166534", bg: "#dcfce7" },
  error:      { text: "エラー",    color: "#991b1b", bg: "#fee2e2" },
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
    <main style={{ maxWidth: 860, margin: "0 auto", padding: "2rem 1.5rem" }}>
      {/* ヘッダー */}
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "1.5rem" }}>
        <div>
          <h1 style={{ margin: 0, fontSize: "1.4rem" }}>投稿した動画</h1>
          <p style={{ margin: "0.25rem 0 0", fontSize: "0.85rem", color: "#888" }}>
            {loading ? "" : `${videos.length} 件`}
          </p>
        </div>
        <Link href="/videos/new" style={{
          padding: "0.5rem 1.2rem", background: "#e00", color: "#fff",
          textDecoration: "none", borderRadius: 6, fontSize: "0.9rem",
        }}>
          + 新規投稿
        </Link>
      </div>

      {/* テーブル */}
      {loading ? (
        <p style={{ color: "#aaa", padding: "2rem 0" }}>読み込み中...</p>
      ) : videos.length === 0 ? (
        <div style={{
          textAlign: "center", padding: "4rem 0", color: "#aaa",
          border: "2px dashed #e5e5e5", borderRadius: 12,
        }}>
          <div style={{ fontSize: "3rem", marginBottom: "1rem" }}>📭</div>
          <p style={{ margin: "0 0 1rem" }}>まだ投稿した動画がありません</p>
          <Link href="/videos/new" style={{ color: "#e00", textDecoration: "none", fontSize: "0.9rem" }}>
            最初の動画を投稿する →
          </Link>
        </div>
      ) : (
        <div style={{ border: "1px solid #e5e5e5", borderRadius: 10, overflow: "hidden", background: "#fff" }}>
          {/* テーブルヘッダー */}
          <div style={{
            display: "grid", gridTemplateColumns: "72px 1fr 110px 120px 100px",
            padding: "0.6rem 1rem", background: "#f9f9f9",
            borderBottom: "1px solid #e5e5e5", fontSize: "0.75rem",
            color: "#888", fontWeight: 600, letterSpacing: "0.05em",
          }}>
            <span></span>
            <span>タイトル</span>
            <span>ステータス</span>
            <span>投稿日</span>
            <span></span>
          </div>

          {/* 行 */}
          {videos.map((video, i) => {
            const st = statusLabel[video.status] ?? statusLabel.done;
            return (
              <div key={video.id} style={{
                display: "grid", gridTemplateColumns: "72px 1fr 110px 120px 100px",
                alignItems: "center", padding: "0.75rem 1rem",
                borderBottom: i < videos.length - 1 ? "1px solid #f0f0f0" : "none",
                transition: "background 0.1s",
              }}
                onMouseEnter={e => (e.currentTarget as HTMLElement).style.background = "#fafafa"}
                onMouseLeave={e => (e.currentTarget as HTMLElement).style.background = "transparent"}
              >
                {/* サムネイル */}
                <Link href={`/videos/${video.id}`} style={{ textDecoration: "none" }}>
                  <div style={{ width: 64, height: 40, background: "#111", borderRadius: 4, overflow: "hidden" }}>
                    {video.thumbnail_url ? (
                      <img src={video.thumbnail_url} alt="" style={{ width: "100%", height: "100%", objectFit: "cover" }} />
                    ) : (
                      <div style={{ height: "100%", display: "flex", alignItems: "center", justifyContent: "center" }}>
                        <span style={{ color: "#555", fontSize: "1rem" }}>▶</span>
                      </div>
                    )}
                  </div>
                </Link>

                {/* タイトル */}
                <Link href={`/videos/${video.id}`} style={{ textDecoration: "none", color: "#111", overflow: "hidden" }}>
                  <p style={{ margin: 0, fontWeight: 600, fontSize: "0.9rem", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                    {video.title}
                  </p>
                </Link>

                {/* ステータスバッジ */}
                <span style={{
                  display: "inline-block", padding: "0.2rem 0.6rem",
                  background: st.bg, color: st.color,
                  borderRadius: 20, fontSize: "0.75rem", fontWeight: 600,
                  width: "fit-content",
                }}>
                  {st.text}
                </span>

                {/* 日付 */}
                <span style={{ fontSize: "0.8rem", color: "#aaa" }}>
                  {new Date(video.created_at).toLocaleDateString("ja-JP")}
                </span>

                {/* 操作 */}
                <div style={{ display: "flex", gap: "0.4rem", justifyContent: "flex-end" }}>
                  <Link href={`/videos/${video.id}/edit`} style={{
                    padding: "0.3rem 0.7rem", border: "1px solid #ddd", borderRadius: 4,
                    textDecoration: "none", color: "#555", fontSize: "0.8rem",
                    background: "#fff",
                  }}>
                    編集
                  </Link>
                  <button onClick={() => handleDelete(video.id)} disabled={deletingId === video.id}
                    style={{
                      padding: "0.3rem 0.7rem", border: "none", borderRadius: 4,
                      background: deletingId === video.id ? "#eee" : "#fee2e2",
                      color: deletingId === video.id ? "#aaa" : "#e00",
                      cursor: deletingId === video.id ? "not-allowed" : "pointer",
                      fontSize: "0.8rem",
                    }}>
                    削除
                  </button>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </main>
  );
}
