"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

type Video = {
  id: number;
  title: string;
  description: string;
  thumbnail_url: string;
  created_at: string;
};

export default function Home() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(true);
  const [query, setQuery] = useState("");
  const [appliedQuery, setAppliedQuery] = useState(""); // 実際に検索に使ったキーワード

  const fetchVideos = (q: string) => {
    setLoading(true);
    setAppliedQuery(q);
    const url = q
      ? `${process.env.NEXT_PUBLIC_API_URL}/api/videos?q=${encodeURIComponent(q)}`
      : `${process.env.NEXT_PUBLIC_API_URL}/api/videos`;
    fetch(url)
      .then((res) => (res.ok ? res.json() : []))
      .then((data) => setVideos(Array.isArray(data) ? data : []))
      .catch(() => setVideos([]))
      .finally(() => setLoading(false));
  };

  useEffect(() => { fetchVideos(""); }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    fetchVideos(query);
  };

  const handleClear = () => {
    setQuery("");
    fetchVideos("");
  };

  return (
    <main style={{ maxWidth: 1100, margin: "0 auto", padding: "2rem 1.5rem" }}>

      {/* 検索バー */}
      <form onSubmit={handleSubmit} style={{ display: "flex", gap: "0.5rem", marginBottom: "1rem" }}>
        <div style={{ position: "relative", flex: 1 }}>
          <span style={{
            position: "absolute", left: "0.85rem", top: "50%", transform: "translateY(-50%)",
            color: "#aaa", fontSize: "1rem", pointerEvents: "none",
          }}>🔍</span>
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="動画を検索..."
            style={{
              width: "100%", padding: "0.65rem 2.2rem 0.65rem 2.4rem",
              fontSize: "0.95rem", border: "1px solid #ddd", borderRadius: 8,
              outline: "none", boxSizing: "border-box", background: "#fff",
            }}
          />
          {query && (
            <button type="button" onClick={handleClear} style={{
              position: "absolute", right: "0.75rem", top: "50%", transform: "translateY(-50%)",
              background: "none", border: "none", color: "#aaa", cursor: "pointer",
              fontSize: "1rem", lineHeight: 1, padding: 0,
            }}>
              ✕
            </button>
          )}
        </div>
        <button type="submit" style={{
          padding: "0 1.25rem", background: "#e00", color: "#fff",
          border: "none", borderRadius: 8, cursor: "pointer",
          fontSize: "0.95rem", fontWeight: 600, whiteSpace: "nowrap",
        }}>
          検索
        </button>
      </form>

      {/* 件数 or キーワード表示 */}
      <p style={{ margin: "0 0 1rem", fontSize: "0.85rem", color: "#888" }}>
        {appliedQuery
          ? `「${appliedQuery}」の検索結果：${loading ? "..." : `${videos.length} 件`}`
          : loading ? "" : `${videos.length} 件`}
      </p>

      {loading ? (
        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))", gap: "1.5rem" }}>
          {[...Array(6)].map((_, i) => (
            <div key={i} style={{ borderRadius: 8, overflow: "hidden" }}>
              <div style={{ background: "#e8e8e8", height: 160, borderRadius: 8 }} />
              <div style={{ padding: "0.75rem 0" }}>
                <div style={{ background: "#e8e8e8", height: 14, borderRadius: 4, marginBottom: 8 }} />
                <div style={{ background: "#e8e8e8", height: 12, borderRadius: 4, width: "60%" }} />
              </div>
            </div>
          ))}
        </div>
      ) : videos.length === 0 ? (
        <div style={{ textAlign: "center", padding: "5rem 0", color: "#aaa" }}>
          <div style={{ fontSize: "3rem", marginBottom: "1rem" }}>
            {appliedQuery ? "🔍" : "🎬"}
          </div>
          <p style={{ fontSize: "1rem", margin: "0 0 1.5rem" }}>
            {appliedQuery ? `「${appliedQuery}」に一致する動画はありません` : "まだ動画がありません"}
          </p>
          {!appliedQuery && (
            <Link href="/videos/new" style={{
              display: "inline-block", padding: "0.75rem 1.5rem",
              background: "#e00", color: "#fff", textDecoration: "none", borderRadius: 6,
            }}>
              最初の動画を投稿する
            </Link>
          )}
        </div>
      ) : (
        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))", gap: "1.5rem" }}>
          {videos.map((video) => (
            <Link key={video.id} href={`/videos/${video.id}`} style={{ textDecoration: "none", color: "inherit" }}>
              <div
                style={{
                  borderRadius: 10, overflow: "hidden", transition: "transform 0.15s, box-shadow 0.15s",
                  cursor: "pointer", border: "1px solid #e5e5e5", background: "#fff",
                }}
                onMouseEnter={e => {
                  (e.currentTarget as HTMLElement).style.transform = "translateY(-4px)";
                  (e.currentTarget as HTMLElement).style.boxShadow = "0 8px 24px rgba(0,0,0,0.12)";
                }}
                onMouseLeave={e => {
                  (e.currentTarget as HTMLElement).style.transform = "translateY(0)";
                  (e.currentTarget as HTMLElement).style.boxShadow = "none";
                }}
              >
                <div style={{ background: "#1a1a1a", height: 165, overflow: "hidden" }}>
                  {video.thumbnail_url ? (
                    <img src={video.thumbnail_url} alt={video.title}
                      style={{ width: "100%", height: "100%", objectFit: "cover" }} />
                  ) : (
                    <div style={{ height: "100%", display: "flex", alignItems: "center", justifyContent: "center" }}>
                      <span style={{ color: "#444", fontSize: "2.5rem" }}>▶</span>
                    </div>
                  )}
                </div>
                <div style={{ padding: "0.75rem 0.75rem 0.75rem" }}>
                  <p style={{
                    margin: "0 0 0.3rem", fontWeight: 600, fontSize: "0.95rem",
                    overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap", color: "#111",
                  }}>
                    {video.title}
                  </p>
                  {video.description && (
                    <p style={{
                      margin: "0 0 0.3rem", fontSize: "0.8rem", color: "#666",
                      overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap",
                    }}>
                      {video.description}
                    </p>
                  )}
                  <p style={{ margin: 0, fontSize: "0.75rem", color: "#aaa" }}>
                    {new Date(video.created_at).toLocaleDateString("ja-JP")}
                  </p>
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </main>
  );
}
