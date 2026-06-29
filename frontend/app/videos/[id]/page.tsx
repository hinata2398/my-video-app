import VideoDetailClient from "./VideoDetailClient";

export async function generateStaticParams() {
  return [{ id: "_" }];
}

export const dynamicParams = false;

export default function Page() {
  return <VideoDetailClient />;
}
