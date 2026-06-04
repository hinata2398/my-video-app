import Link from "next/link";

type Video = {
  id: number;
  title: string;
  description: string;
  thumbnail_url: string;
  created_at: string;
};

async function getVideos(): Promise<Video[]> {
  try {
    const url = process.env.INTERNAL_API_URL ?? process.env.NEXT_PUBLIC_API_URL;
    const res = await fetch(`${url}/api/videos`, { cache: "no-store" });
    if (!res.ok) return [];
    const data = await res.json();
    return Array.isArray(data) ? data : [];
  } catch {
    return [];
  }
}

export default async function Home() {
  const videos = await getVideos();

  return (
    <main style={{ maxWidth: 900, margin: "0 auto", padding: "2rem", fontFamily: "sans-serif" }}>
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "2rem" }}>
        <h1 style={{ margin: 0 }}>動画一覧</h1>
        <Link href="/videos/new" style={{ padding: "0.5rem 1rem", background: "#e00", color: "#fff", textDecoration: "none", borderRadius: 4 }}>
          + 投稿する
        </Link>
      </div>
      {videos.length === 0 ? (
        <div style={{ textAlign: "center", padding: "4rem 0", color: "#888" }}>
          <div style={{ fontSize: "4rem", marginBottom: "1rem" }}>▶</div>
          <p style={{ fontSize: "1.1rem", margin: "0 0 1rem" }}>まだ動画がありません</p>
          <Link href="/videos/new" style={{ color: "#e00", textDecoration: "none" }}>
            最初の動画を投稿する →
          </Link>
        </div>
      ) : (
        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(260px, 1fr))", gap: "1.5rem" }}>
          {videos.map((video) => (
            <Link key={video.id} href={`/videos/${video.id}`} style={{ textDecoration: "none", color: "inherit" }}>
              <div style={{ border: "1px solid #ddd", borderRadius: 8, overflow: "hidden" }}>
                <div style={{ background: "#111", height: 150, display: "flex", alignItems: "center", justifyContent: "center" }}>
                  {video.thumbnail_url ? (
                    <img src={video.thumbnail_url} alt={video.title} style={{ width: "100%", height: "100%", objectFit: "cover" }} />
                  ) : (
                    <span style={{ color: "#555", fontSize: "2rem" }}>▶</span>
                  )}
                </div>
                <div style={{ padding: "0.75rem" }}>
                  <p style={{ margin: 0, fontWeight: "bold", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{video.title}</p>
                  <p style={{ margin: "0.25rem 0 0", fontSize: "0.8rem", color: "#888" }}>
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
