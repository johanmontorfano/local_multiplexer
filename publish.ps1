$comment = $args[0]

Write-Output "[publish-ps] Publishing a new commit using Git."
Write-Output "[publish-ps] Comment: $comment"

git.exe add .
git.exe commit -am $comment
git.exe push
git.exe fetch