require 'bundler'
Bundler.setup()
require 'pry'

namespace :xdr do

  # As diamcircle-core adds more .x files, we'll need to update this array
  # Prior to launch, we should be separating our .x files into a separate
  # repo, and should be able to improve this integration.
  HAYASHI_XDR = [
                  "src/xdr/diamcircle-SCP.x",
                  "src/xdr/diamcircle-ledger-entries.x",
                  "src/xdr/diamcircle-ledger.x",
                  "src/xdr/diamcircle-overlay.x",
                  "src/xdr/diamcircle-transaction.x",
                  "src/xdr/diamcircle-types.x"
                ]
  LOCAL_XDR_PATHS = HAYASHI_XDR.map{ |src| "xdr/" + File.basename(src) }

  task :update => [:download, :generate]

  task :download do
    require 'octokit'
    require 'base64'
    FileUtils.mkdir_p "xdr"
    FileUtils.rm_rf "xdr/*.x"

    client = Octokit::Client.new(:netrc => true)

    HAYASHI_XDR.each do |src|
      local_path = "xdr/" + File.basename(src)
      encoded    = client.contents("diamcircle/diamcircle-core", path: src).content
      decoded    = Base64.decode64 encoded

      IO.write(local_path, decoded)
    end
  end

  task :generate do
    require "pathname"
    require "xdrgen"
    require 'fileutils'

    compilation = Xdrgen::Compilation.new(
      LOCAL_XDR_PATHS,
      output_dir: "xdr",
      namespace:  "xdr",
      language:   :go
    )
    compilation.compile

    xdr_generated = IO.read("xdr/xdr_generated.go")
    IO.write("xdr/xdr_generated.go", <<~EOS)
      //lint:file-ignore S1005 The issue should be fixed in xdrgen. Unfortunately, there's no way to ignore a single file in staticcheck.
      //lint:file-ignore U1000 fmtTest is not needed anywhere, should be removed in xdrgen.
      #{xdr_generated}
    EOS

    system("gofmt -w xdr/xdr_generated.go")
  end
end
