function handler(event) {
  var req = event.request;
  var uri = req.uri;

  // 1) 動的: /videos/<数値>/edit → 編集シェル
  if (/^\/videos\/\d+\/edit\/?$/.test(uri)) {
    req.uri = "/videos/_/edit.html";
    return req;
  }
  // 2) 動的: /videos/<数値> → 詳細シェル
  if (/^\/videos\/\d+\/?$/.test(uri)) {
    req.uri = "/videos/_.html";
    return req;
  }
  // 3) 静的クリーンURL: 拡張子が無く "/" でもない → .html を付ける
  //    （/_next/...js などは "." を含むので対象外、/ は default_root_object に任せる）
  if (uri !== "/" && !uri.includes(".")) {
    req.uri = uri.replace(/\/$/, "") + ".html";
  }
  return req;
}
