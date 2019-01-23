param(
    [int]$DebPackageVersion = 1
)

function ConvertFrom-SecureStringPlain {
    param(
        [Parameter(ValueFromPipeline=$true,Mandatory=$true,Position=0)]
        [System.Security.SecureString]
        $sstr
    )
    $marshal = [System.Runtime.InteropServices.Marshal]
    $ptr = $marshal::SecureStringToBSTR( $sstr )
    $str = $marshal::PtrToStringBSTR( $ptr )
    $marshal::ZeroFreeBSTR( $ptr )
    $str
}

function New-FileObject {
    [cmdletBinding()]
    param(
        [string]$repo,
        [string]$version,
        [string]$packageVersion,
        [string]$arch
    )
    return [PSCustomObject]@{
        repo = $repo
        version = $version
        packageVersion = $packageVersion
        arch = $arch
    }
}

function Get-SourceFilename{
    param(
        [PSCustomObject]$FileObject
    )
    return "loxwebhook_v$($FileObject.version)-1_$($FileObject.arch).deb"
}

function Get-TargetFilename{
    param(
        [PSCustomObject]$FileObject
    )
    return "loxwebhook_v$($FileObject.version)-$($FileObject.packageVersion)_$($FileObject.arch).deb"
}

function Get-URI{
    param(
        [PSCustomObject]$FileObject
    )
    return "https://api.bintray.com/content/axel/"+
                 "$($FileObject.repo)/loxwebhook/"+
                 "$($FileObject.version)-$($FileObject.packageVersion)/"+
                 "$(Get-TargetFilename($FileObject))"+
                 ";deb_distribution=stretch"+
                 ";deb_component=main"+
                 ";deb_architecture=$($FileObject.arch)"+
                 ";publish=1"
}

& .\set_deploy_env.ps1
$files = @()
$branch = &git rev-parse --abbrev-ref HEAD
switch ($branch) {
    "master" {
        $debRepo = "loxwebhook_deb"
        $version = (&git describe --tags --abbrev=0).substring(1)
        $proc_goreleaser = Start-Process -FilePath 'goreleaser.exe' -ArgumentList "--rm-dist" -NoNewWindow -Wait -ErrorAction Stop -PassThru
        $files += ((New-FileObject -repo $debRepo -version $version -packageVersion $DebPackageVersion -arch "arm"))
        $files += ((New-FileObject -repo $debRepo -version $version -packageVersion $DebPackageVersion -arch "amd64"))
    }
    "dev" {
        $debRepo = "loxwebhook_deb_dev"
        $currentCommit = &git rev-parse --short HEAD
        $version = (&git describe --tags --abbrev=0).substring(1) + ".$currentCommit"
        $proc_goreleaser = Start-Process -FilePath 'goreleaser.exe' -ArgumentList "--rm-dist", "--snapshot" -NoNewWindow -Wait -ErrorAction Stop -PassThru
        $files += ((New-FileObject -repo $debRepo -version $version -packageVersion $DebPackageVersion -arch "arm"))
        $files += ((New-FileObject -repo $debRepo -version $version -packageVersion $DebPackageVersion -arch "amd64"))
    }
    default {
        Write-Error "Wrong branch: $branch" -ErrorAction Stop
    }
}
if ($proc_goreleaser.ExitCode -gt 0) {
    Exit $?
}

$pw = ConvertTo-SecureString -String $env:BINTRAY_API_KEY -AsPlainText -Force -ErrorAction Stop
$user = $env:BINTRAY_USERNAME
$cred = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $user, $pw -ErrorAction Stop
$gpg_key_pw = Read-Host -Prompt 'Password for gpg signing key' -AsSecureString

$headers = @{"X-GPG-PASSPHRASE" = ConvertFrom-SecureStringPlain($gpg_key_pw)}
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

foreach ($f in $files) {
    if ($f.packageVersion -ne 1) {
        Rename-Item -Path "./dist/$(Get-SourceFilename -FileObject $f)" -NewName (Get-TargetFilename -FileObject $f)
    }
    Write-Host (Get-URI -FileObject $f)
    try {
        Invoke-RestMethod -Uri (Get-URI -FileObject $f) -Method Put -InFile "./dist/$(Get-TargetFilename -FileObject $f)" -Headers $headers -Credential $cred
    }
    catch {
        $response = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($response)
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Error $_
        Write-Output $responseBody
        Exit(1)
    }
}
