class AppstoreScraper < Formula
  desc "JSON-first CLI for Apple App Store and Google Play public data"
  homepage "https://github.com/youngminz/appstore-scraper-cli"
  url "https://github.com/youngminz/appstore-scraper-cli/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_RELEASE_TARBALL_SHA256"
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
