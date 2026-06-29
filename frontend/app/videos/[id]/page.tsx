import VideoDetailClient from "./VideoDetailClient";

export async function generateStaticParams() {
  try {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
    const res = await fetch(`${apiUrl}/api/videos`);
    const videos = await res.json();
    if (videos.length > 0) {
      return videos.map((v: { id: number }) => ({ id: String(v.id) }));
    }
  } catch {
    // API unavailable at build time
  }
  return [{ id: "_" }];
}

export const dynamicParams = false;

export default function Page() {
  return <VideoDetailClient />;
}
