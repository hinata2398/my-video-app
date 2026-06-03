import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "My Video App",
  description: "動画投稿サイト ポートフォリオ",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
