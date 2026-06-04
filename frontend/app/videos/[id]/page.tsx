import Link from "next/link";

type Video = {
  id: number;
  user_id: number;
  title: string;
  description: string;
  thumbnail_url: string;
  created_at: string;
};

async function getVideo(id: string): Promise<Video | null> {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/videos/${id}`, {
    cache: "no-store",
  });
  if (!res.ok) return null;
  return res.json();
}

export default async function VideoDetailPage({ params }: { params: { id: string } }) {
  const video = await getVideo(params.id);

  if (!video) {
    return (
      <main style={{ maxWidth: 800, margin: "0 auto", padding: "2rem", fontFamily: "sans-serif" }}>
        <p>動画が見つかりません。</p>
        <Link href="/">← 一覧へ戻る</Link>
      </main>
    );
  }

  return (
    <main style={{ maxWidth: 800, margin: "0 auto", padding: "2rem", fontFamily: "sans-serif" }}>
      <Link href="/" style={{ color: "#888", textDecoration: "none" }}>← 一覧へ戻る</Link>
      <div style={{ background: "#111", marginTop: "1rem", borderRadius: 8, height: 400, display: "flex", alignItems: "center", justifyContent: "center" }}>
        {video.thumbnail_url ? (
          <img src={video.thumbnail_url} alt={video.title} style={{ width: "100%", height: "100%", objectFit: "cover", borderRadius: 8 }} />
        ) : (
          <span style={{ color: "#555", fontSize: "4rem" }}>▶</span>
        )}
      </div>
      <h1 style={{ marginTop: "1rem" }}>{video.title}</h1>
      <p style={{ color: "#888", fontSize: "0.9rem" }}>{new Date(video.created_at).toLocaleDateString("ja-JP")}</p>
      {video.description && <p style={{ marginTop: "1rem", lineHeight: 1.7 }}>{video.description}</p>}
    </main>
  );
}
