#define MyAppName "FluxStream"
#define MyAppVersion "2.0.0"
#define MyAppPublisher "FluxStream"
#define MyAppExeName "fluxstream.exe"
#ifndef SourceDir
  #define SourceDir "C:\xampp\htdocs\stream\dist\fluxstream-windows-amd64-service"
#endif
#ifndef OutputDir
  #define OutputDir "C:\xampp\htdocs\stream\dist"
#endif
#ifndef SetupIconPath
  #define SetupIconPath "C:\xampp\htdocs\stream\deployment\fluxstream.ico"
#endif

[Setup]
AppId={{8A7F7A82-8D71-4BA1-90A0-E34D2BFEA5C7}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppPublisher={#MyAppPublisher}
DefaultDirName={autopf}\FluxStream
DefaultGroupName=FluxStream
PrivilegesRequired=admin
OutputDir={#OutputDir}
OutputBaseFilename=FluxStream-Setup
Compression=lzma
SolidCompression=yes
WizardStyle=modern
ArchitecturesInstallIn64BitMode=x64compatible
DisableProgramGroupPage=yes
SetupIconFile={#SetupIconPath}
UninstallDisplayIcon={app}\fluxstream.ico

[Languages]
Name: "turkish"; MessagesFile: "compiler:Languages\Turkish.isl"

[Tasks]
Name: "desktopicon"; Description: "Masaustu kisayolu"; Flags: checkedonce

[Files]
Source: "{#SourceDir}\fluxstream.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#SourceDir}\fluxstream.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#SourceDir}\ffmpeg\*"; DestDir: "{app}\ffmpeg"; Flags: ignoreversion recursesubdirs createallsubdirs

[UninstallDelete]
Type: filesandordirs; Name: "{app}\data"
Type: filesandordirs; Name: "{app}\ffmpeg"
Type: files; Name: "{app}\fluxstream.exe"
Type: files; Name: "{app}\fluxstream.ico"
Type: dirifempty; Name: "{app}"

[Icons]
Name: "{group}\FluxStream"; Filename: "{app}\fluxstream.exe"; IconFilename: "{app}\fluxstream.ico"
Name: "{autodesktop}\FluxStream"; Filename: "{app}\fluxstream.exe"; Tasks: desktopicon; IconFilename: "{app}\fluxstream.ico"

[Run]
Filename: "{app}\fluxstream.exe"; Description: "FluxStream yonetim paneli icin uygulamayi baslat"; Flags: nowait postinstall skipifsilent; Check: ShouldLaunchDesktopMode

[UninstallRun]
Filename: "{app}\fluxstream.exe"; Parameters: "service stop"; Flags: runhidden waituntilterminated skipifdoesntexist; RunOnceId: "FluxStreamServiceStop"
Filename: "{app}\fluxstream.exe"; Parameters: "service uninstall"; Flags: runhidden waituntilterminated skipifdoesntexist; RunOnceId: "FluxStreamServiceUninstall"

[Code]
var
  ModePage: TInputOptionWizardPage;
  NetworkPage: TWizardPage;
  SSLPage: TWizardPage;
  UseCustomPortsCheck: TNewCheckBox;
  HttpPortLabel: TNewStaticText;
  HttpPortEdit: TNewEdit;
  RtmpPortLabel: TNewStaticText;
  RtmpPortEdit: TNewEdit;
  HttpsPortLabel: TNewStaticText;
  HttpsPortEdit: TNewEdit;
  RtmpsPortLabel: TNewStaticText;
  RtmpsPortEdit: TNewEdit;
  DomainLabel: TNewStaticText;
  DomainEdit: TNewEdit;
  PreloadSSLCheck: TNewCheckBox;
  EnableSSLNowCheck: TNewCheckBox;
  CertPathLabel: TNewStaticText;
  CertPathEdit: TNewEdit;
  CertBrowseButton: TNewButton;
  KeyPathLabel: TNewStaticText;
  KeyPathEdit: TNewEdit;
  KeyBrowseButton: TNewButton;

function BoolToConfig(Value: Boolean): String;
begin
  if Value then
    Result := 'true'
  else
    Result := 'false';
end;

procedure AddConfigArg(var Params: String; const Key, Value: String);
begin
  if Params <> '' then
    Params := Params + ' ';
  Params := Params + '"' + Key + '=' + Value + '"';
end;

function ShouldInstallService(): Boolean;
begin
  Result := ModePage.Values[0];
end;

function ShouldLaunchDesktopMode(): Boolean;
begin
  Result := not ShouldInstallService();
end;

procedure UpdatePortPageState(Sender: TObject);
begin
  HttpPortEdit.Enabled := UseCustomPortsCheck.Checked;
  RtmpPortEdit.Enabled := UseCustomPortsCheck.Checked;
  HttpsPortEdit.Enabled := UseCustomPortsCheck.Checked;
  RtmpsPortEdit.Enabled := UseCustomPortsCheck.Checked;
  HttpPortLabel.Enabled := UseCustomPortsCheck.Checked;
  RtmpPortLabel.Enabled := UseCustomPortsCheck.Checked;
  HttpsPortLabel.Enabled := UseCustomPortsCheck.Checked;
  RtmpsPortLabel.Enabled := UseCustomPortsCheck.Checked;
end;

procedure UpdateSSLPageState(Sender: TObject);
begin
  CertPathEdit.Enabled := PreloadSSLCheck.Checked;
  CertBrowseButton.Enabled := PreloadSSLCheck.Checked;
  KeyPathEdit.Enabled := PreloadSSLCheck.Checked;
  KeyBrowseButton.Enabled := PreloadSSLCheck.Checked;
  CertPathLabel.Enabled := PreloadSSLCheck.Checked;
  KeyPathLabel.Enabled := PreloadSSLCheck.Checked;
  EnableSSLNowCheck.Enabled := PreloadSSLCheck.Checked;
end;

procedure BrowseForCert(Sender: TObject);
var
  SelectedFile: String;
begin
  SelectedFile := CertPathEdit.Text;
  if GetOpenFileName('CRT dosyasi secin', SelectedFile, '', 'Certificate Files|*.crt;*.pem|Tum Dosyalar|*.*', '') then
    CertPathEdit.Text := SelectedFile;
end;

procedure BrowseForKey(Sender: TObject);
var
  SelectedFile: String;
begin
  SelectedFile := KeyPathEdit.Text;
  if GetOpenFileName('KEY dosyasi secin', SelectedFile, '', 'Key Files|*.key;*.pem|Tum Dosyalar|*.*', '') then
    KeyPathEdit.Text := SelectedFile;
end;

function ValidatePortValue(const Caption, Value: String): Boolean;
var
  Port: Integer;
begin
  Port := StrToIntDef(Value, -1);
  Result := (Port > 0) and (Port <= 65535);
  if not Result then
    MsgBox(Caption + ' icin 1-65535 araliginda gecerli bir port girin.', mbError, MB_OK);
end;

function ServiceExists(): Boolean;
var
  ResultCode: Integer;
begin
  Result := Exec(ExpandConstant('{sys}\sc.exe'), 'query FluxStream', '', SW_HIDE, ewWaitUntilTerminated, ResultCode) and (ResultCode = 0);
end;

function RunInstalledCommand(const Params, ErrorMessage: String): Boolean;
var
  ResultCode: Integer;
begin
  Result := Exec(ExpandConstant('{app}\fluxstream.exe'), Params, '', SW_HIDE, ewWaitUntilTerminated, ResultCode) and (ResultCode = 0);
  if not Result then
    RaiseException(ErrorMessage);
end;

procedure ApplyInstallerSettings();
var
  Params: String;
  CertTarget: String;
  KeyTarget: String;
begin
  Params := '';

  if UseCustomPortsCheck.Checked then
  begin
    AddConfigArg(Params, 'http_port', HttpPortEdit.Text);
    AddConfigArg(Params, 'embed_http_port', HttpPortEdit.Text);
    AddConfigArg(Params, 'rtmp_port', RtmpPortEdit.Text);
    AddConfigArg(Params, 'https_port', HttpsPortEdit.Text);
    AddConfigArg(Params, 'embed_https_port', HttpsPortEdit.Text);
    AddConfigArg(Params, 'rtmps_port', RtmpsPortEdit.Text);
  end;

  if Trim(DomainEdit.Text) <> '' then
    AddConfigArg(Params, 'embed_domain', Trim(DomainEdit.Text));

  if PreloadSSLCheck.Checked then
  begin
    ForceDirectories(ExpandConstant('{app}\data\certs'));
    CertTarget := ExpandConstant('{app}\data\certs\server.crt');
    KeyTarget := ExpandConstant('{app}\data\certs\server.key');
    if not CopyFile(CertPathEdit.Text, CertTarget, False) then
      RaiseException('SSL sertifika dosyasi kopyalanamadi.');
    if not CopyFile(KeyPathEdit.Text, KeyTarget, False) then
      RaiseException('SSL key dosyasi kopyalanamadi.');

    AddConfigArg(Params, 'ssl_cert_path', CertTarget);
    AddConfigArg(Params, 'ssl_key_path', KeyTarget);

    if EnableSSLNowCheck.Checked then
    begin
      AddConfigArg(Params, 'ssl_enabled', BoolToConfig(True));
      AddConfigArg(Params, 'rtmps_enabled', BoolToConfig(True));
      AddConfigArg(Params, 'embed_use_https', BoolToConfig(True));
    end;
  end;

  if Params <> '' then
    RunInstalledCommand('config set ' + Params, 'Kurulum secenekleri uygulanamadi.');
end;

procedure ConfigureInstallMode();
var
  ResultCode: Integer;
begin
  if ServiceExists() then
  begin
    Exec(ExpandConstant('{app}\fluxstream.exe'), 'service stop', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
    Exec(ExpandConstant('{app}\fluxstream.exe'), 'service uninstall', '', SW_HIDE, ewWaitUntilTerminated, ResultCode);
  end;

  if ShouldInstallService() then
  begin
    RunInstalledCommand('service install', 'FluxStream servisi kurulamadi.');
    RunInstalledCommand('service start', 'FluxStream servisi baslatilamadi.');
  end;
end;

procedure InitializeWizard();
var
  NoteLabel: TNewStaticText;
  CurrentTop: Integer;
begin
  ModePage := CreateInputOptionPage(
    wpSelectTasks,
    'Kurulum Modu',
    'FluxStream''i nasil calistirmak istediginizi secin',
    'Bu secim istege baglidir. Kurulumdan sonra admin panelinden, Windows service araclarindan veya yeniden kurulumla degistirebilirsiniz.',
    True,
    False);
  ModePage.Add('Windows Service olarak kur ve otomatik baslat (onerilen)');
  ModePage.Add('Servis kurma, uygulamayi elle/konsoldan baslat');
  ModePage.Values[0] := True;

  NetworkPage := CreateCustomPage(
    ModePage.ID,
    'Ag Ayarlari (Opsiyonel)',
    'Isterseniz varsayilan portlari kurulum sirasinda ozellestirebilirsiniz. Bu ayarlar kurulumdan sonra admin panelinden degistirilebilir.');

  NoteLabel := TNewStaticText.Create(NetworkPage);
  NoteLabel.Parent := NetworkPage.Surface;
  NoteLabel.Left := 0;
  NoteLabel.Top := 0;
  NoteLabel.Width := NetworkPage.SurfaceWidth;
  NoteLabel.Height := ScaleY(42);
  NoteLabel.WordWrap := True;
  NoteLabel.Caption := 'Bu sayfadaki ayarlar istege baglidir. Domain ve web portlari embed/player linklerinde kullanilir. RTMP ve RTMPS portlari yayin gonderimi icin ayridir. Tum ayarlar kurulumdan sonra panelden degistirilebilir.';

  UseCustomPortsCheck := TNewCheckBox.Create(NetworkPage);
  UseCustomPortsCheck.Parent := NetworkPage.Surface;
  UseCustomPortsCheck.Left := 0;
  UseCustomPortsCheck.Top := NoteLabel.Top + NoteLabel.Height + ScaleY(8);
  UseCustomPortsCheck.Width := NetworkPage.SurfaceWidth;
  UseCustomPortsCheck.Caption := 'Kurulum sirasinda ozel portlari uygula';
  UseCustomPortsCheck.OnClick := @UpdatePortPageState;

  CurrentTop := UseCustomPortsCheck.Top + ScaleY(32);

  DomainLabel := TNewStaticText.Create(NetworkPage);
  DomainLabel.Parent := NetworkPage.Surface;
  DomainLabel.Left := 0;
  DomainLabel.Top := CurrentTop;
  DomainLabel.Caption := 'Public Domain / IP';

  DomainEdit := TNewEdit.Create(NetworkPage);
  DomainEdit.Parent := NetworkPage.Surface;
  DomainEdit.Left := ScaleX(180);
  DomainEdit.Top := CurrentTop - ScaleY(4);
  DomainEdit.Width := ScaleX(260);
  DomainEdit.Text := '';

  CurrentTop := CurrentTop + ScaleY(32);

  HttpPortLabel := TNewStaticText.Create(NetworkPage);
  HttpPortLabel.Parent := NetworkPage.Surface;
  HttpPortLabel.Left := 0;
  HttpPortLabel.Top := CurrentTop;
  HttpPortLabel.Caption := 'HTTP Port';

  HttpPortEdit := TNewEdit.Create(NetworkPage);
  HttpPortEdit.Parent := NetworkPage.Surface;
  HttpPortEdit.Left := ScaleX(180);
  HttpPortEdit.Top := CurrentTop - ScaleY(4);
  HttpPortEdit.Width := ScaleX(120);
  HttpPortEdit.Text := '8844';

  CurrentTop := CurrentTop + ScaleY(28);
  RtmpPortLabel := TNewStaticText.Create(NetworkPage);
  RtmpPortLabel.Parent := NetworkPage.Surface;
  RtmpPortLabel.Left := 0;
  RtmpPortLabel.Top := CurrentTop;
  RtmpPortLabel.Caption := 'RTMP Port';

  RtmpPortEdit := TNewEdit.Create(NetworkPage);
  RtmpPortEdit.Parent := NetworkPage.Surface;
  RtmpPortEdit.Left := ScaleX(180);
  RtmpPortEdit.Top := CurrentTop - ScaleY(4);
  RtmpPortEdit.Width := ScaleX(120);
  RtmpPortEdit.Text := '1935';

  CurrentTop := CurrentTop + ScaleY(28);
  HttpsPortLabel := TNewStaticText.Create(NetworkPage);
  HttpsPortLabel.Parent := NetworkPage.Surface;
  HttpsPortLabel.Left := 0;
  HttpsPortLabel.Top := CurrentTop;
  HttpsPortLabel.Caption := 'HTTPS Port';

  HttpsPortEdit := TNewEdit.Create(NetworkPage);
  HttpsPortEdit.Parent := NetworkPage.Surface;
  HttpsPortEdit.Left := ScaleX(180);
  HttpsPortEdit.Top := CurrentTop - ScaleY(4);
  HttpsPortEdit.Width := ScaleX(120);
  HttpsPortEdit.Text := '443';

  CurrentTop := CurrentTop + ScaleY(28);
  RtmpsPortLabel := TNewStaticText.Create(NetworkPage);
  RtmpsPortLabel.Parent := NetworkPage.Surface;
  RtmpsPortLabel.Left := 0;
  RtmpsPortLabel.Top := CurrentTop;
  RtmpsPortLabel.Caption := 'RTMPS Port';

  RtmpsPortEdit := TNewEdit.Create(NetworkPage);
  RtmpsPortEdit.Parent := NetworkPage.Surface;
  RtmpsPortEdit.Left := ScaleX(180);
  RtmpsPortEdit.Top := CurrentTop - ScaleY(4);
  RtmpsPortEdit.Width := ScaleX(120);
  RtmpsPortEdit.Text := '1936';

  SSLPage := CreateCustomPage(
    NetworkPage.ID,
    'SSL On Yukleme (Opsiyonel)',
    'CRT ve KEY dosyalarini simdiden ekleyebilirsiniz. Istemezseniz kurulumdan sonra admin panelinden yukleyebilirsiniz.');

  NoteLabel := TNewStaticText.Create(SSLPage);
  NoteLabel.Parent := SSLPage.Surface;
  NoteLabel.Left := 0;
  NoteLabel.Top := 0;
  NoteLabel.Width := SSLPage.SurfaceWidth;
  NoteLabel.Height := ScaleY(42);
  NoteLabel.WordWrap := True;
  NoteLabel.Caption := 'SSL adimi istege baglidir. Dosya secerseniz sertifikalar uygulamanin data/certs klasorune kopyalanir. Tum ayarlar kurulumdan sonra admin panelinden degistirilebilir.';

  PreloadSSLCheck := TNewCheckBox.Create(SSLPage);
  PreloadSSLCheck.Parent := SSLPage.Surface;
  PreloadSSLCheck.Left := 0;
  PreloadSSLCheck.Top := NoteLabel.Top + NoteLabel.Height + ScaleY(8);
  PreloadSSLCheck.Width := SSLPage.SurfaceWidth;
  PreloadSSLCheck.Caption := 'Kurulum sirasinda CRT ve KEY dosyalarini kopyala';
  PreloadSSLCheck.OnClick := @UpdateSSLPageState;

  EnableSSLNowCheck := TNewCheckBox.Create(SSLPage);
  EnableSSLNowCheck.Parent := SSLPage.Surface;
  EnableSSLNowCheck.Left := ScaleX(16);
  EnableSSLNowCheck.Top := PreloadSSLCheck.Top + ScaleY(24);
  EnableSSLNowCheck.Width := SSLPage.SurfaceWidth - ScaleX(16);
  EnableSSLNowCheck.Caption := 'Kurulumdan sonra HTTPS ve RTMPS''i aktif et';
  EnableSSLNowCheck.Checked := True;

  CurrentTop := EnableSSLNowCheck.Top + ScaleY(32);

  CertPathLabel := TNewStaticText.Create(SSLPage);
  CertPathLabel.Parent := SSLPage.Surface;
  CertPathLabel.Left := 0;
  CertPathLabel.Top := CurrentTop;
  CertPathLabel.Caption := 'CRT Dosyasi';

  CertPathEdit := TNewEdit.Create(SSLPage);
  CertPathEdit.Parent := SSLPage.Surface;
  CertPathEdit.Left := 0;
  CertPathEdit.Top := CurrentTop + ScaleY(16);
  CertPathEdit.Width := SSLPage.SurfaceWidth - ScaleX(96);

  CertBrowseButton := TNewButton.Create(SSLPage);
  CertBrowseButton.Parent := SSLPage.Surface;
  CertBrowseButton.Left := CertPathEdit.Left + CertPathEdit.Width + ScaleX(8);
  CertBrowseButton.Top := CertPathEdit.Top - ScaleY(1);
  CertBrowseButton.Width := ScaleX(80);
  CertBrowseButton.Caption := 'Sec...';
  CertBrowseButton.OnClick := @BrowseForCert;

  CurrentTop := CertPathEdit.Top + ScaleY(34);
  KeyPathLabel := TNewStaticText.Create(SSLPage);
  KeyPathLabel.Parent := SSLPage.Surface;
  KeyPathLabel.Left := 0;
  KeyPathLabel.Top := CurrentTop;
  KeyPathLabel.Caption := 'KEY Dosyasi';

  KeyPathEdit := TNewEdit.Create(SSLPage);
  KeyPathEdit.Parent := SSLPage.Surface;
  KeyPathEdit.Left := 0;
  KeyPathEdit.Top := CurrentTop + ScaleY(16);
  KeyPathEdit.Width := SSLPage.SurfaceWidth - ScaleX(96);

  KeyBrowseButton := TNewButton.Create(SSLPage);
  KeyBrowseButton.Parent := SSLPage.Surface;
  KeyBrowseButton.Left := KeyPathEdit.Left + KeyPathEdit.Width + ScaleX(8);
  KeyBrowseButton.Top := KeyPathEdit.Top - ScaleY(1);
  KeyBrowseButton.Width := ScaleX(80);
  KeyBrowseButton.Caption := 'Sec...';
  KeyBrowseButton.OnClick := @BrowseForKey;

  UpdatePortPageState(nil);
  UpdateSSLPageState(nil);
end;

function NextButtonClick(CurPageID: Integer): Boolean;
begin
  Result := True;

  if CurPageID = NetworkPage.ID then
  begin
    if UseCustomPortsCheck.Checked then
      Result :=
        ValidatePortValue('HTTP Port', HttpPortEdit.Text) and
        ValidatePortValue('RTMP Port', RtmpPortEdit.Text) and
        ValidatePortValue('HTTPS Port', HttpsPortEdit.Text) and
        ValidatePortValue('RTMPS Port', RtmpsPortEdit.Text);
  end;

  if Result and (CurPageID = SSLPage.ID) and PreloadSSLCheck.Checked then
  begin
    if (Trim(CertPathEdit.Text) = '') or (not FileExists(CertPathEdit.Text)) then
    begin
      MsgBox('Gecerli bir CRT dosyasi secin veya SSL on yuklemeyi kapatin.', mbError, MB_OK);
      Result := False;
    end
    else if (Trim(KeyPathEdit.Text) = '') or (not FileExists(KeyPathEdit.Text)) then
    begin
      MsgBox('Gecerli bir KEY dosyasi secin veya SSL on yuklemeyi kapatin.', mbError, MB_OK);
      Result := False;
    end;
  end;
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    ApplyInstallerSettings();
    ConfigureInstallMode();
  end;
end;

function InitializeSetup(): Boolean;
begin
  Result := True;
end;
