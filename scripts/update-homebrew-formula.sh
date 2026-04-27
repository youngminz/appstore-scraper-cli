#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 3 ]]; then
  echo "usage: $0 <version-tag> <sha256> <formula-path>" >&2
  exit 2
fi

version_tag="$1"
sha256="$2"
formula_path="$3"

cat >"${formula_path}" <<EOF
class AppstoreScraper < Formula
  desc "JSON-first CLI for Apple App Store and Google Play public data"
  homepage "https://github.com/youngminz/appstore-scraper-cli"
  url "https://github.com/youngminz/appstore-scraper-cli/archive/refs/tags/${version_tag}.tar.gz"
  sha256 "${sha256}"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w")
  end

  test do
    assert_match "appstore-scraper retrieves public mobile app store data",
      shell_output("#{bin}/appstore-scraper --help")
  end
end
EOF
