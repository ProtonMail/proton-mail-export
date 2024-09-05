#include "ExportedData.h"
#include <fstream>
#include <iostream>

//****************************************************************************************************************************************************
/// \brief write content into a file.
/// \param[in] path The  path of the folder to write int.
/// \param[in] content The content to write into the file
//****************************************************************************************************************************************************
void writeFile(std::filesystem::path const& path, std::string const& content) {
    std::ofstream stream(path, std::ios::out);
    stream << content;
}

//****************************************************************************************************************************************************
/// Once we decide to switch to C++23, we can replace this with binary resource inclusion: see https://en.cppreference.com/w/c/preprocessor/embed
///
/// \param[in] dir The folder path. The folder must exist.
//****************************************************************************************************************************************************
void createTestBackup(std::filesystem::path const& dir) {
    writeFile(dir / "4HGNtgH-oj3CPMiABEMYt9se-38NAM8b4T7h2JfIq0D3WGsTa0tZxFhzNmnoFApweAbSUS1rbm9mxspaC6MOxQ==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <4HGNtgH-oj3CPMiABEMYt9se-38NAM8b4T7h2JfIq0D3WGsTa0tZxFhzNmnoFApweAbSUS1rbm9mxspaC6MOxQ==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:03:47 +0000
X-Pm-External-Id: <tMDbpBIpnbtFHTaFuNKqCI2rDCt-aR0YEZJ9HvRmFCnMykWFY4B4zYHn5JkElg1Kc1Yz87fBU8_f9QXvwv0XhEcx5g77KLr9uzdYWQZMIKw=@michelon.ch>
X-Pm-Internal-Id: 4HGNtgH-oj3CPMiABEMYt9se-38NAM8b4T7h2JfIq0D3WGsTa0tZxFhzNmnoFApweAbSUS1rbm9mxspaC6MOxQ==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: Act 3 - Scene 3
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:03:48 +0000
X-Pm-Spamscore: 0
Message-Id: <tMDbpBIpnbtFHTaFuNKqCI2rDCt-aR0YEZJ9HvRmFCnMykWFY4B4zYHn5JkElg1Kc1Yz87fBU8_f9QXvwv0XhEcx5g77KLr9uzdYWQZMIKw=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:03:47 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>DESDE=
MONA</span><div><span>Be thou assured, good Cassio, I will do</span></div><=
span>All my abilities in thy behalf</span>.<br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>
)");

    writeFile(dir / "4HGNtgH-oj3CPMiABEMYt9se-38NAM8b4T7h2JfIq0D3WGsTa0tZxFhzNmnoFApweAbSUS1rbm9mxspaC6MOxQ==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "4HGNtgH-oj3CPMiABEMYt9se-38NAM8b4T7h2JfIq0D3WGsTa0tZxFhzNmnoFApweAbSUS1rbm9mxspaC6MOxQ==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "0",
      "5",
      "15",
      "VDJ9fXR7xsoyWzBeU-Nr3tWT-0hc0dRkQWKnr2kp8l8gzSKq13oWLaGKdFDDNtRx17TkJqtns_vDS-BGWL3fEA=="
    ],
    "ExternalID": "tMDbpBIpnbtFHTaFuNKqCI2rDCt-aR0YEZJ9HvRmFCnMykWFY4B4zYHn5JkElg1Kc1Yz87fBU8_f9QXvwv0XhEcx5g77KLr9uzdYWQZMIKw=@michelon.ch",
    "Subject": "Act 3 - Scene 3",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718345027,
    "Size": 774,
    "Unread": 1,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 3 - Scene 3\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:03:47 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003ctMDbpBIpnbtFHTaFuNKqCI2rDCt-aR0YEZJ9HvRmFCnMykWFY4B4zYHn5JkElg1Kc1Yz87fBU8_f9QXvwv0XhEcx5g77KLr9uzdYWQZMIKw=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:03:48 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "8m8cT5h2ZwEX2y_iGfSOOWWFaGEYk8NlRLKblc9X_aWxbCC6sjyknJBILbe9m7yUJhixXHePSqHD_4Zqx59fBA==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <8m8cT5h2ZwEX2y_iGfSOOWWFaGEYk8NlRLKblc9X_aWxbCC6sjyknJBILbe9m7yUJhixXHePSqHD_4Zqx59fBA==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:01:20 +0000
X-Pm-External-Id: <oyaVbqwyeZiYh52oJDpRSBPX_qkqIixmSesvsT0DRM-RqBPM0AyYHVL8wQ8m7-EnjX3u-7RXJnTKLOsH2_tCxY-Jqsyc_5vt8VHCX-rw_x8=@michelon.ch>
X-Pm-Internal-Id: 8m8cT5h2ZwEX2y_iGfSOOWWFaGEYk8NlRLKblc9X_aWxbCC6sjyknJBILbe9m7yUJhixXHePSqHD_4Zqx59fBA==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: Act 2 - Scene 2
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:01:25 +0000
X-Pm-Spamscore: 0
Message-Id: <oyaVbqwyeZiYh52oJDpRSBPX_qkqIixmSesvsT0DRM-RqBPM0AyYHVL8wQ8m7-EnjX3u-7RXJnTKLOsH2_tCxY-Jqsyc_5vt8VHCX-rw_x8=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:01:20 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>HERAL=
D &nbsp;It is Othello's pleasure, our noble and valiant</span><div><span>ge=
neral, that upon certain tidings now arrived,</span></div><div><span>import=
ing the mere perdition of the Turkish fleet,</span></div><div><span>every m=
an put himself into triumph: some to</span></div><div><span>dance, some to =
make bonfires, each man to what</span></div><span>sport and revels his addi=
tion leads him.</span><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>
)");

    writeFile(dir / "8m8cT5h2ZwEX2y_iGfSOOWWFaGEYk8NlRLKblc9X_aWxbCC6sjyknJBILbe9m7yUJhixXHePSqHD_4Zqx59fBA==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "8m8cT5h2ZwEX2y_iGfSOOWWFaGEYk8NlRLKblc9X_aWxbCC6sjyknJBILbe9m7yUJhixXHePSqHD_4Zqx59fBA==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "5",
      "15",
      "UD2CCpFokJimuXRGM91VglWoj4eOLCevCs0JIfVDcmlLmlqUxVcqWW_6UMyz7D17XBbcqK4SDxUBO66KY8LlCw==",
      "wKYBxlcMwhPUHACLkhISvx462lDsU0PSpdZF3p9ki9viXb1kIam5KEByamy9CFSl6HPTVXIEIBt1hwMpJBmXOQ=="
    ],
    "ExternalID": "oyaVbqwyeZiYh52oJDpRSBPX_qkqIixmSesvsT0DRM-RqBPM0AyYHVL8wQ8m7-EnjX3u-7RXJnTKLOsH2_tCxY-Jqsyc_5vt8VHCX-rw_x8=@michelon.ch",
    "Subject": "Act 2 - Scene 2",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718344880,
    "Size": 1053,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 2 - Scene 2\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:01:20 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003coyaVbqwyeZiYh52oJDpRSBPX_qkqIixmSesvsT0DRM-RqBPM0AyYHVL8wQ8m7-EnjX3u-7RXJnTKLOsH2_tCxY-Jqsyc_5vt8VHCX-rw_x8=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:01:25 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "B0BdjEJQwG3nUoys4u16JlvyGUtvJeuc6ZRqItsvvyxQmc5-N3fuFoSGMRPsE9_VG2O3zU-XWbG9uiAWnUMsPQ==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <B0BdjEJQwG3nUoys4u16JlvyGUtvJeuc6ZRqItsvvyxQmc5-N3fuFoSGMRPsE9_VG2O3zU-XWbG9uiAWnUMsPQ==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 05:59:10 +0000
X-Pm-External-Id: <50tTfSemGK4A6Np_GYucejnpyEMgvxeN-R1oy7JfqCuFggVdIuuwOZyvt3GLcilOXydWmtFBgV5ZPKyjssiC85I_0zC5mSGpvcMmeH2QDVc=@michelon.ch>
X-Pm-Internal-Id: B0BdjEJQwG3nUoys4u16JlvyGUtvJeuc6ZRqItsvvyxQmc5-N3fuFoSGMRPsE9_VG2O3zU-XWbG9uiAWnUMsPQ==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: scene 3
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 05:59:14 +0000
X-Pm-Spamscore: 0
Message-Id: <50tTfSemGK4A6Np_GYucejnpyEMgvxeN-R1oy7JfqCuFggVdIuuwOZyvt3GLcilOXydWmtFBgV5ZPKyjssiC85I_0zC5mSGpvcMmeH2QDVc=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 05:59:10 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>DUKE,=
 [reading a paper]</span><div><span>There's no composition in these news</s=
pan></div><div><span>That gives them credit.</span></div><div><br></div><di=
v><span>FIRST SENATOR, [reading a paper]</span></div><div><span>Indeed, the=
y are disproportioned.</span></div><span>My letters say a hundred and seven=
 galleys.</span><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>
)");

    writeFile(dir / "B0BdjEJQwG3nUoys4u16JlvyGUtvJeuc6ZRqItsvvyxQmc5-N3fuFoSGMRPsE9_VG2O3zU-XWbG9uiAWnUMsPQ==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "B0BdjEJQwG3nUoys4u16JlvyGUtvJeuc6ZRqItsvvyxQmc5-N3fuFoSGMRPsE9_VG2O3zU-XWbG9uiAWnUMsPQ==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "5",
      "15",
      "UD2CCpFokJimuXRGM91VglWoj4eOLCevCs0JIfVDcmlLmlqUxVcqWW_6UMyz7D17XBbcqK4SDxUBO66KY8LlCw==",
      "bW35iHyu46lo65YIHCkG7qW2IpEw87PXDFZQ8Zx2JoZ5BLQAz4h4JuKOhPEbF-JEZMwDjVc5clSiOE37KEVj_A=="
    ],
    "ExternalID": "50tTfSemGK4A6Np_GYucejnpyEMgvxeN-R1oy7JfqCuFggVdIuuwOZyvt3GLcilOXydWmtFBgV5ZPKyjssiC85I_0zC5mSGpvcMmeH2QDVc=@michelon.ch",
    "Subject": "scene 3",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718344750,
    "Size": 972,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: scene 3\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 05:59:10 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003c50tTfSemGK4A6Np_GYucejnpyEMgvxeN-R1oy7JfqCuFggVdIuuwOZyvt3GLcilOXydWmtFBgV5ZPKyjssiC85I_0zC5mSGpvcMmeH2QDVc=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 05:59:14 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "drTldxNF7Ae4qNz23sw-jl8rheRZk0OSKKFh0IoQLSUQcwjM92kcxsezHbR2gSMQl0pHXuz-lLf_ueFPm22LRw==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <drTldxNF7Ae4qNz23sw-jl8rheRZk0OSKKFh0IoQLSUQcwjM92kcxsezHbR2gSMQl0pHXuz-lLf_ueFPm22LRw==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:01:53 +0000
X-Pm-External-Id: <vJgUJK_Uw4DimMkaasz9lcr2hNHg9w6JQOublDXbD7KL3crINy1o-aCkJyXUMJqOp8Nuoii1YJgByjK5F2Kd8COiXoFTA0D4qV3mpD_xjvg=@michelon.ch>
X-Pm-Internal-Id: drTldxNF7Ae4qNz23sw-jl8rheRZk0OSKKFh0IoQLSUQcwjM92kcxsezHbR2gSMQl0pHXuz-lLf_ueFPm22LRw==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: Act 2 - Scene 3
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:01:57 +0000
X-Pm-Spamscore: 0
Message-Id: <vJgUJK_Uw4DimMkaasz9lcr2hNHg9w6JQOublDXbD7KL3crINy1o-aCkJyXUMJqOp8Nuoii1YJgByjK5F2Kd8COiXoFTA0D4qV3mpD_xjvg=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:01:53 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>OTHEL=
LO</span><div><span>Good Michael, look you to the guard tonight.</span></di=
v><div><span>Let's teach ourselves that honorable stop</span></div><span>No=
t to outsport discretion.</span><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>
)");

    writeFile(dir / "drTldxNF7Ae4qNz23sw-jl8rheRZk0OSKKFh0IoQLSUQcwjM92kcxsezHbR2gSMQl0pHXuz-lLf_ueFPm22LRw==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "drTldxNF7Ae4qNz23sw-jl8rheRZk0OSKKFh0IoQLSUQcwjM92kcxsezHbR2gSMQl0pHXuz-lLf_ueFPm22LRw==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "5",
      "15",
      "UD2CCpFokJimuXRGM91VglWoj4eOLCevCs0JIfVDcmlLmlqUxVcqWW_6UMyz7D17XBbcqK4SDxUBO66KY8LlCw==",
      "wKYBxlcMwhPUHACLkhISvx462lDsU0PSpdZF3p9ki9viXb1kIam5KEByamy9CFSl6HPTVXIEIBt1hwMpJBmXOQ=="
    ],
    "ExternalID": "vJgUJK_Uw4DimMkaasz9lcr2hNHg9w6JQOublDXbD7KL3crINy1o-aCkJyXUMJqOp8Nuoii1YJgByjK5F2Kd8COiXoFTA0D4qV3mpD_xjvg=@michelon.ch",
    "Subject": "Act 2 - Scene 3",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718344913,
    "Size": 838,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 2 - Scene 3\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:01:53 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nMessage-Id: \u003cvJgUJK_Uw4DimMkaasz9lcr2hNHg9w6JQOublDXbD7KL3crINy1o-aCkJyXUMJqOp8Nuoii1YJgByjK5F2Kd8COiXoFTA0D4qV3mpD_xjvg=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:01:57 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "GUzJVBVHWUHHB3B6M5aii8SrYOPReEv-ab3SOzn1wDvYru7TCaVtDje5VUGT0t4GLcpKzFiiDUpKUJy_6r8erg==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <GUzJVBVHWUHHB3B6M5aii8SrYOPReEv-ab3SOzn1wDvYru7TCaVtDje5VUGT0t4GLcpKzFiiDUpKUJy_6r8erg==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:04:23 +0000
X-Pm-External-Id: <SHQKSObUPz7hxp0zmxsb56KXxEQ4Mf-nwpi6A1hlGJmigKTyQi4xNLqxqt3D8Z4F0KuUKwXOXDOfhgHHGK2Wra9Nm2i1VcxfEnbYTh6guyY=@michelon.ch>
X-Pm-Internal-Id: GUzJVBVHWUHHB3B6M5aii8SrYOPReEv-ab3SOzn1wDvYru7TCaVtDje5VUGT0t4GLcpKzFiiDUpKUJy_6r8erg==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: Act 3 - Scene 4
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:04:26 +0000
X-Pm-Spamscore: 0
Message-Id: <SHQKSObUPz7hxp0zmxsb56KXxEQ4Mf-nwpi6A1hlGJmigKTyQi4xNLqxqt3D8Z4F0KuUKwXOXDOfhgHHGK2Wra9Nm2i1VcxfEnbYTh6guyY=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:04:23 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>DESDE=
MONA &nbsp;Do you know, sirrah, where Lieutenant</span><br><span>Cassio lie=
s?</span><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>
)");

    writeFile(dir / "GUzJVBVHWUHHB3B6M5aii8SrYOPReEv-ab3SOzn1wDvYru7TCaVtDje5VUGT0t4GLcpKzFiiDUpKUJy_6r8erg==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "GUzJVBVHWUHHB3B6M5aii8SrYOPReEv-ab3SOzn1wDvYru7TCaVtDje5VUGT0t4GLcpKzFiiDUpKUJy_6r8erg==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "0",
      "5",
      "15",
      "VDJ9fXR7xsoyWzBeU-Nr3tWT-0hc0dRkQWKnr2kp8l8gzSKq13oWLaGKdFDDNtRx17TkJqtns_vDS-BGWL3fEA=="
    ],
    "ExternalID": "SHQKSObUPz7hxp0zmxsb56KXxEQ4Mf-nwpi6A1hlGJmigKTyQi4xNLqxqt3D8Z4F0KuUKwXOXDOfhgHHGK2Wra9Nm2i1VcxfEnbYTh6guyY=@michelon.ch",
    "Subject": "Act 3 - Scene 4",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718345063,
    "Size": 740,
    "Unread": 1,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 3 - Scene 4\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:04:23 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003cSHQKSObUPz7hxp0zmxsb56KXxEQ4Mf-nwpi6A1hlGJmigKTyQi4xNLqxqt3D8Z4F0KuUKwXOXDOfhgHHGK2Wra9Nm2i1VcxfEnbYTh6guyY=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:04:26 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "KihWtsTJuZbzK1HXRVYJdMMFq2sdsSqSb29WLVYiub7WwnqLWDEPQbiXExgnNJJIOo7zesVa20nNtJIW3S-lHA==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <KihWtsTJuZbzK1HXRVYJdMMFq2sdsSqSb29WLVYiub7WwnqLWDEPQbiXExgnNJJIOo7zesVa20nNtJIW3S-lHA==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:03:21 +0000
X-Pm-External-Id: <EN_g7aKnrGfPwjq8UrqJSGFOEjBdvoxujOvXrapuDK81IrdtEGsIJbecvV0oxpPRHNXLnzo4oLk5vIrMpyhWqmTdPP3BGtGt_bu0Dpgdr8c=@michelon.ch>
X-Pm-Internal-Id: KihWtsTJuZbzK1HXRVYJdMMFq2sdsSqSb29WLVYiub7WwnqLWDEPQbiXExgnNJJIOo7zesVa20nNtJIW3S-lHA==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: Act 3 - Scene 2
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:03:23 +0000
X-Pm-Spamscore: 0
Message-Id: <EN_g7aKnrGfPwjq8UrqJSGFOEjBdvoxujOvXrapuDK81IrdtEGsIJbecvV0oxpPRHNXLnzo4oLk5vIrMpyhWqmTdPP3BGtGt_bu0Dpgdr8c=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:03:21 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>CASSI=
O</span><div><span>Masters, play here (I will content your pains)</span></d=
iv><div><span>Something that's brief; and bid "Good morrow,</span></div><sp=
an>general."	[They play.]</span><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>
)");

    writeFile(dir / "KihWtsTJuZbzK1HXRVYJdMMFq2sdsSqSb29WLVYiub7WwnqLWDEPQbiXExgnNJJIOo7zesVa20nNtJIW3S-lHA==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "KihWtsTJuZbzK1HXRVYJdMMFq2sdsSqSb29WLVYiub7WwnqLWDEPQbiXExgnNJJIOo7zesVa20nNtJIW3S-lHA==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "0",
      "5",
      "15",
      "VDJ9fXR7xsoyWzBeU-Nr3tWT-0hc0dRkQWKnr2kp8l8gzSKq13oWLaGKdFDDNtRx17TkJqtns_vDS-BGWL3fEA=="
    ],
    "ExternalID": "EN_g7aKnrGfPwjq8UrqJSGFOEjBdvoxujOvXrapuDK81IrdtEGsIJbecvV0oxpPRHNXLnzo4oLk5vIrMpyhWqmTdPP3BGtGt_bu0Dpgdr8c=@michelon.ch",
    "Subject": "Act 3 - Scene 2",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718345001,
    "Size": 838,
    "Unread": 1,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 3 - Scene 2\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:03:21 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003cEN_g7aKnrGfPwjq8UrqJSGFOEjBdvoxujOvXrapuDK81IrdtEGsIJbecvV0oxpPRHNXLnzo4oLk5vIrMpyhWqmTdPP3BGtGt_bu0Dpgdr8c=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:03:23 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "KOE0LUr_ZoIafHj76-buUaOO694Er4eWrkmfKHkbWipVJOsZD7urz_oAaIS0KnkxEcbip9lVKIaYW4gYGzNUvQ==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <KOE0LUr_ZoIafHj76-buUaOO694Er4eWrkmfKHkbWipVJOsZD7urz_oAaIS0KnkxEcbip9lVKIaYW4gYGzNUvQ==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:04:43 +0000
X-Pm-External-Id: <gx3Vcq5qhDHi6MMjab65ld-z5JjgywSe0Jt44nKfKXywC72IcoVQHRNtm5D38y_XMIS9qrowlLQ42yULuBtCMGivMxIivQgQNW3xu1oaAX0=@michelon.ch>
X-Pm-Internal-Id: KOE0LUr_ZoIafHj76-buUaOO694Er4eWrkmfKHkbWipVJOsZD7urz_oAaIS0KnkxEcbip9lVKIaYW4gYGzNUvQ==
To: "dev@michelon.ch" <dev@michelon.ch>
Reply-To: "Test @ Michelon" <test@michelon.ch>
From: "Test @ Michelon" <test@michelon.ch>
Subject: Act 4 - Scene 1
X-Pm-Recipient-Encryption: dev%40michelon.ch=pgp-pm
X-Pm-Recipient-Authentication: dev%40michelon.ch=pgp-pm
X-Pm-Scheduled-Sent-Original-Time: Fri, 14 Jun 2024 06:04:29 +0000
Message-Id: <gx3Vcq5qhDHi6MMjab65ld-z5JjgywSe0Jt44nKfKXywC72IcoVQHRNtm5D38y_XMIS9qrowlLQ42yULuBtCMGivMxIivQgQNW3xu1oaAX0=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:04:43 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>IAGO<=
/span><br><span>Will you think so?</span><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">

            </div>

            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">

            </div>
</div>
)");

    writeFile(dir / "KOE0LUr_ZoIafHj76-buUaOO694Er4eWrkmfKHkbWipVJOsZD7urz_oAaIS0KnkxEcbip9lVKIaYW4gYGzNUvQ==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "KOE0LUr_ZoIafHj76-buUaOO694Er4eWrkmfKHkbWipVJOsZD7urz_oAaIS0KnkxEcbip9lVKIaYW4gYGzNUvQ==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "2",
      "5",
      "7",
      "15"
    ],
    "ExternalID": "gx3Vcq5qhDHi6MMjab65ld-z5JjgywSe0Jt44nKfKXywC72IcoVQHRNtm5D38y_XMIS9qrowlLQ42yULuBtCMGivMxIivQgQNW3xu1oaAX0=@michelon.ch",
    "Subject": "Act 4 - Scene 1",
    "Sender": {
      "Name": "Test @ Michelon",
      "Address": "test@michelon.ch"
    },
    "ToList": [
      {
        "Name": "dev@michelon.ch",
        "Address": "dev@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "Flags": 8206,
    "Time": 1718345083,
    "Size": 677,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 4 - Scene 1\r\nTo: dev@michelon.ch \u003cdev@michelon.ch\u003e\r\nFrom: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:04:43 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003cgx3Vcq5qhDHi6MMjab65ld-z5JjgywSe0Jt44nKfKXywC72IcoVQHRNtm5D38y_XMIS9qrowlLQ42yULuBtCMGivMxIivQgQNW3xu1oaAX0=@michelon.ch\u003e\r\nX-Pm-Scheduled-Sent-Original-Time: Fri, 14 Jun 2024 06:04:29 +0000\r\nX-Pm-Recipient-Authentication: dev%40michelon.ch=pgp-pm\r\nX-Pm-Recipient-Encryption: dev%40michelon.ch=pgp-pm\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "L1GdzPQiz6cC4dlJSiTN-nVc9gvbVpC67ShGvmvje8qwI-QkcrOWbQtJgd0tBOFa2ERxWoCuOnThSxcZIwV4Zw==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <L1GdzPQiz6cC4dlJSiTN-nVc9gvbVpC67ShGvmvje8qwI-QkcrOWbQtJgd0tBOFa2ERxWoCuOnThSxcZIwV4Zw==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 05:57:50 +0000
X-Pm-External-Id: <2UiYOZ_7W_POO5g9sv5QMjk3eSNri3JGibSbrAY89jbS9csDCvs1g5agtsNccetayBtM0PD2SwGIK4r5tDyInmjNAKdEYI-fNbJwMfTAPlU=@michelon.ch>
X-Pm-Internal-Id: L1GdzPQiz6cC4dlJSiTN-nVc9gvbVpC67ShGvmvje8qwI-QkcrOWbQtJgd0tBOFa2ERxWoCuOnThSxcZIwV4Zw==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: Scene 1
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 05:57:55 +0000
X-Pm-Spamscore: 0
Message-Id: <2UiYOZ_7W_POO5g9sv5QMjk3eSNri3JGibSbrAY89jbS9csDCvs1g5agtsNccetayBtM0PD2SwGIK4r5tDyInmjNAKdEYI-fNbJwMfTAPlU=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 05:57:50 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>RODER=
IGO</span><div><span>Tush, never tell me! I take it much unkindly</span></d=
iv><div><span>That thou, Iago, who hast had my purse</span></div><div><span=
>As if the strings were thine, shouldst know of this.</span></div><div><br>=
</div><div><span>IAGO &nbsp;'Sblood, but you'll not hear me!</span></div><d=
iv><span>If ever I did dream of such a matter,</span></div><span>Abhor me.<=
/span><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>
)");

    writeFile(dir / "L1GdzPQiz6cC4dlJSiTN-nVc9gvbVpC67ShGvmvje8qwI-QkcrOWbQtJgd0tBOFa2ERxWoCuOnThSxcZIwV4Zw==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "L1GdzPQiz6cC4dlJSiTN-nVc9gvbVpC67ShGvmvje8qwI-QkcrOWbQtJgd0tBOFa2ERxWoCuOnThSxcZIwV4Zw==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "5",
      "15",
      "UD2CCpFokJimuXRGM91VglWoj4eOLCevCs0JIfVDcmlLmlqUxVcqWW_6UMyz7D17XBbcqK4SDxUBO66KY8LlCw==",
      "bW35iHyu46lo65YIHCkG7qW2IpEw87PXDFZQ8Zx2JoZ5BLQAz4h4JuKOhPEbF-JEZMwDjVc5clSiOE37KEVj_A=="
    ],
    "ExternalID": "2UiYOZ_7W_POO5g9sv5QMjk3eSNri3JGibSbrAY89jbS9csDCvs1g5agtsNccetayBtM0PD2SwGIK4r5tDyInmjNAKdEYI-fNbJwMfTAPlU=@michelon.ch",
    "Subject": "Scene 1",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718344670,
    "Size": 1037,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Scene 1\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 05:57:50 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003c2UiYOZ_7W_POO5g9sv5QMjk3eSNri3JGibSbrAY89jbS9csDCvs1g5agtsNccetayBtM0PD2SwGIK4r5tDyInmjNAKdEYI-fNbJwMfTAPlU=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 05:57:55 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "labels.json",
              R"({
  "Version": 1,
  "Payload": [
    {
      "Path": "",
      "ID": "0",
      "ParentID": "",
      "Name": "Inbox",
      "Color": "#8080FF",
      "Type": 3
    },
    {
      "Path": "",
      "ID": "8",
      "ParentID": "",
      "Name": "Drafts",
      "Color": "#8080FF",
      "Type": 3
    },
    {
      "Path": "All Drafts",
      "ID": "1",
      "ParentID": "",
      "Name": "All Drafts",
      "Color": "#8080FF",
      "Type": 1
    },
    {
      "Path": "All Scheduled",
      "ID": "12",
      "ParentID": "",
      "Name": "All Scheduled",
      "Color": "#8080FF",
      "Type": 1
    },
    {
      "Path": "",
      "ID": "7",
      "ParentID": "",
      "Name": "Sent",
      "Color": "#8080FF",
      "Type": 3
    },
    {
      "Path": "All Sent",
      "ID": "2",
      "ParentID": "",
      "Name": "All Sent",
      "Color": "#8080FF",
      "Type": 1
    },
    {
      "Path": "Starred",
      "ID": "10",
      "ParentID": "",
      "Name": "Starred",
      "Color": "#8080FF",
      "Type": 1
    },
    {
      "Path": "Snoozed",
      "ID": "16",
      "ParentID": "",
      "Name": "Snoozed",
      "Color": "#8080FF",
      "Type": 1
    },
    {
      "Path": "",
      "ID": "6",
      "ParentID": "",
      "Name": "Archive",
      "Color": "#8080FF",
      "Type": 3
    },
    {
      "Path": "",
      "ID": "4",
      "ParentID": "",
      "Name": "Spam",
      "Color": "#8080FF",
      "Type": 3
    },
    {
      "Path": "",
      "ID": "3",
      "ParentID": "",
      "Name": "Trash",
      "Color": "#8080FF",
      "Type": 3
    },
    {
      "Path": "All Mail",
      "ID": "5",
      "ParentID": "",
      "Name": "All Mail",
      "Color": "#8080FF",
      "Type": 1
    },
    {
      "Path": "All Mail",
      "ID": "15",
      "ParentID": "",
      "Name": "All Mail",
      "Color": "#8080FF",
      "Type": 1
    },
    {
      "Path": "Outbox",
      "ID": "9",
      "ParentID": "",
      "Name": "Outbox",
      "Color": "#8080FF",
      "Type": 1
    },
    {
      "Path": "Othello",
      "ID": "UD2CCpFokJimuXRGM91VglWoj4eOLCevCs0JIfVDcmlLmlqUxVcqWW_6UMyz7D17XBbcqK4SDxUBO66KY8LlCw==",
      "ParentID": "",
      "Name": "Othello",
      "Color": "#DB60D6",
      "Type": 3
    },
    {
      "Path": "Act 1",
      "ID": "bW35iHyu46lo65YIHCkG7qW2IpEw87PXDFZQ8Zx2JoZ5BLQAz4h4JuKOhPEbF-JEZMwDjVc5clSiOE37KEVj_A==",
      "ParentID": "",
      "Name": "Act 1",
      "Color": "#5252cc",
      "Type": 1
    },
    {
      "Path": "Act 2",
      "ID": "wKYBxlcMwhPUHACLkhISvx462lDsU0PSpdZF3p9ki9viXb1kIam5KEByamy9CFSl6HPTVXIEIBt1hwMpJBmXOQ==",
      "ParentID": "",
      "Name": "Act 2",
      "Color": "#3CBB3A",
      "Type": 1
    },
    {
      "Path": "Act 3",
      "ID": "VDJ9fXR7xsoyWzBeU-Nr3tWT-0hc0dRkQWKnr2kp8l8gzSKq13oWLaGKdFDDNtRx17TkJqtns_vDS-BGWL3fEA==",
      "ParentID": "",
      "Name": "Act 3",
      "Color": "#5252cc",
      "Type": 1
    }
  ]
})");

    writeFile(dir / "mUOvICBvCtjcgJwB3xOJpjH8LfxfykIik7g2H-_KJ3MmOlKe011YlFucPSFljfbwF0pJRysajiISyayQNSaylA==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <mUOvICBvCtjcgJwB3xOJpjH8LfxfykIik7g2H-_KJ3MmOlKe011YlFucPSFljfbwF0pJRysajiISyayQNSaylA==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:05:13 +0000
X-Pm-External-Id: <clRFqej2LpbwJg_tiYx102Ptse5u-vrsGCjtHdLsPbBKD-Sfd9r8aVgoy5YLygSt0LyajpaQMG2FQfvKK9w86N7hdj5Z9jnTMMN_c4PU1I4=@michelon.ch>
X-Pm-Internal-Id: mUOvICBvCtjcgJwB3xOJpjH8LfxfykIik7g2H-_KJ3MmOlKe011YlFucPSFljfbwF0pJRysajiISyayQNSaylA==
To: "dev@michelon.ch" <dev@michelon.ch>
Reply-To: "Test @ Michelon" <test@michelon.ch>
From: "Test @ Michelon" <test@michelon.ch>
Subject: Act 4 - Scene 2
X-Pm-Recipient-Encryption: dev%40michelon.ch=pgp-pm
X-Pm-Recipient-Authentication: dev%40michelon.ch=pgp-pm
X-Pm-Scheduled-Sent-Original-Time: Fri, 14 Jun 2024 06:04:59 +0000
Message-Id: <clRFqej2LpbwJg_tiYx102Ptse5u-vrsGCjtHdLsPbBKD-Sfd9r8aVgoy5YLygSt0LyajpaQMG2FQfvKK9w86N7hdj5Z9jnTMMN_c4PU1I4=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:05:13 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>OTHEL=
LO &nbsp;You have seen nothing then?</span><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">

            </div>

            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">

            </div>
</div>
)");

    writeFile(dir / "mUOvICBvCtjcgJwB3xOJpjH8LfxfykIik7g2H-_KJ3MmOlKe011YlFucPSFljfbwF0pJRysajiISyayQNSaylA==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "mUOvICBvCtjcgJwB3xOJpjH8LfxfykIik7g2H-_KJ3MmOlKe011YlFucPSFljfbwF0pJRysajiISyayQNSaylA==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "2",
      "5",
      "7",
      "15"
    ],
    "ExternalID": "clRFqej2LpbwJg_tiYx102Ptse5u-vrsGCjtHdLsPbBKD-Sfd9r8aVgoy5YLygSt0LyajpaQMG2FQfvKK9w86N7hdj5Z9jnTMMN_c4PU1I4=@michelon.ch",
    "Subject": "Act 4 - Scene 2",
    "Sender": {
      "Name": "Test @ Michelon",
      "Address": "test@michelon.ch"
    },
    "ToList": [
      {
        "Name": "dev@michelon.ch",
        "Address": "dev@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "Flags": 8206,
    "Time": 1718345113,
    "Size": 679,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 4 - Scene 2\r\nTo: dev@michelon.ch \u003cdev@michelon.ch\u003e\r\nFrom: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:05:13 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003cclRFqej2LpbwJg_tiYx102Ptse5u-vrsGCjtHdLsPbBKD-Sfd9r8aVgoy5YLygSt0LyajpaQMG2FQfvKK9w86N7hdj5Z9jnTMMN_c4PU1I4=@michelon.ch\u003e\r\nX-Pm-Scheduled-Sent-Original-Time: Fri, 14 Jun 2024 06:04:59 +0000\r\nX-Pm-Recipient-Authentication: dev%40michelon.ch=pgp-pm\r\nX-Pm-Recipient-Encryption: dev%40michelon.ch=pgp-pm\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "tL3HtkSDhIURwqiSRxvAPUzE1_Hfor-vSncsBoC0bq3mAh11K7YO7qvBO7w2vRDglxuC7dP1HuxEjku5qu26IA==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <tL3HtkSDhIURwqiSRxvAPUzE1_Hfor-vSncsBoC0bq3mAh11K7YO7qvBO7w2vRDglxuC7dP1HuxEjku5qu26IA==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 05:58:38 +0000
X-Pm-External-Id: <Y01uYzel6isvDH0SQ0zC1fi_jStYcbnSMAT9eVMm2MSoafFFk6vtMLs8pAvkVuFRKmacWuvgLDD8D6PUcT-HovgTsWNAuW2TNtXXXt4WD-4=@michelon.ch>
X-Pm-Internal-Id: tL3HtkSDhIURwqiSRxvAPUzE1_Hfor-vSncsBoC0bq3mAh11K7YO7qvBO7w2vRDglxuC7dP1HuxEjku5qu26IA==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: Scene 2
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 05:58:41 +0000
X-Pm-Spamscore: 0
Message-Id: <Y01uYzel6isvDH0SQ0zC1fi_jStYcbnSMAT9eVMm2MSoafFFk6vtMLs8pAvkVuFRKmacWuvgLDD8D6PUcT-HovgTsWNAuW2TNtXXXt4WD-4=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 05:58:38 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>IAGO<=
/span><div><span>Though in the trade of war I have slain men,</span></div><=
div><span>Yet do I hold it very stuff o' th' conscience</span></div><div><s=
pan>To do no contrived murder. I lack iniquity</span></div><div><span>Somet=
imes to do me service. Nine or ten times</span></div><div><span>I had thoug=
ht t' have yerked him here under the</span></div><div><span>ribs.</span></d=
iv><div><br></div><div><span>OTHELLO</span></div><div><span>'Tis better as =
it is.</span></div><span></span><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>
)");

    writeFile(dir / "tL3HtkSDhIURwqiSRxvAPUzE1_Hfor-vSncsBoC0bq3mAh11K7YO7qvBO7w2vRDglxuC7dP1HuxEjku5qu26IA==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "tL3HtkSDhIURwqiSRxvAPUzE1_Hfor-vSncsBoC0bq3mAh11K7YO7qvBO7w2vRDglxuC7dP1HuxEjku5qu26IA==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "5",
      "15",
      "UD2CCpFokJimuXRGM91VglWoj4eOLCevCs0JIfVDcmlLmlqUxVcqWW_6UMyz7D17XBbcqK4SDxUBO66KY8LlCw==",
      "bW35iHyu46lo65YIHCkG7qW2IpEw87PXDFZQ8Zx2JoZ5BLQAz4h4JuKOhPEbF-JEZMwDjVc5clSiOE37KEVj_A=="
    ],
    "ExternalID": "Y01uYzel6isvDH0SQ0zC1fi_jStYcbnSMAT9eVMm2MSoafFFk6vtMLs8pAvkVuFRKmacWuvgLDD8D6PUcT-HovgTsWNAuW2TNtXXXt4WD-4=@michelon.ch",
    "Subject": "Scene 2",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718344718,
    "Size": 1138,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Scene 2\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 05:58:38 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003cY01uYzel6isvDH0SQ0zC1fi_jStYcbnSMAT9eVMm2MSoafFFk6vtMLs8pAvkVuFRKmacWuvgLDD8D6PUcT-HovgTsWNAuW2TNtXXXt4WD-4=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 05:58:41 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "UWIv7DnWeMxRhhPFiMssgVZ3MEnJ3adTR8ZVLm7gTrj8SQJTt9P2wNcKx6q3mTahy9cSgNkyl82nXJRKLHKgKw==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <UWIv7DnWeMxRhhPFiMssgVZ3MEnJ3adTR8ZVLm7gTrj8SQJTt9P2wNcKx6q3mTahy9cSgNkyl82nXJRKLHKgKw==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:00:06 +0000
X-Pm-External-Id: <dXpK0xplB4X-yW7rb40v0WEgzwFd7IZqXD0OoaAENn5hb_my-c9jspdor8gIQsbxGkYA35CsNM9YzR6KxtTnjhL_6hDdFLvHk-9jiNi2ADs=@michelon.ch>
X-Pm-Internal-Id: UWIv7DnWeMxRhhPFiMssgVZ3MEnJ3adTR8ZVLm7gTrj8SQJTt9P2wNcKx6q3mTahy9cSgNkyl82nXJRKLHKgKw==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: Act 2
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:00:09 +0000
X-Pm-Spamscore: 0
Message-Id: <dXpK0xplB4X-yW7rb40v0WEgzwFd7IZqXD0OoaAENn5hb_my-c9jspdor8gIQsbxGkYA35CsNM9YzR6KxtTnjhL_6hDdFLvHk-9jiNi2ADs=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:00:06 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>MONTA=
NO</span><div><span>What from the cape can you discern at sea?</span></div>=
<div><br></div><div><span>FIRST GENTLEMAN</span></div><div><span>Nothing at=
 all. It is a high-wrought flood.</span></div><div><span>I cannot 'twixt th=
e heaven and the main</span></div><div><span>Descry a sail.</span></div><di=
v><br></div></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>
)");

    writeFile(dir / "UWIv7DnWeMxRhhPFiMssgVZ3MEnJ3adTR8ZVLm7gTrj8SQJTt9P2wNcKx6q3mTahy9cSgNkyl82nXJRKLHKgKw==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "UWIv7DnWeMxRhhPFiMssgVZ3MEnJ3adTR8ZVLm7gTrj8SQJTt9P2wNcKx6q3mTahy9cSgNkyl82nXJRKLHKgKw==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "5",
      "15",
      "UD2CCpFokJimuXRGM91VglWoj4eOLCevCs0JIfVDcmlLmlqUxVcqWW_6UMyz7D17XBbcqK4SDxUBO66KY8LlCw==",
      "wKYBxlcMwhPUHACLkhISvx462lDsU0PSpdZF3p9ki9viXb1kIam5KEByamy9CFSl6HPTVXIEIBt1hwMpJBmXOQ=="
    ],
    "ExternalID": "dXpK0xplB4X-yW7rb40v0WEgzwFd7IZqXD0OoaAENn5hb_my-c9jspdor8gIQsbxGkYA35CsNM9YzR6KxtTnjhL_6hDdFLvHk-9jiNi2ADs=@michelon.ch",
    "Subject": "Act 2",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718344806,
    "Size": 964,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 2\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:00:06 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003cdXpK0xplB4X-yW7rb40v0WEgzwFd7IZqXD0OoaAENn5hb_my-c9jspdor8gIQsbxGkYA35CsNM9YzR6KxtTnjhL_6hDdFLvHk-9jiNi2ADs=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:00:09 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "Ve78Yi7YkjntEfYQ4OfWCzJMKmJBCbw-Ybip7P4E531hP1ttzQRq_OuXITt8Z5BXTXxHNJ0p_h6jYPmvfi6B7w==.eml",
              R"(Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8
References: <Ve78Yi7YkjntEfYQ4OfWCzJMKmJBCbw-Ybip7P4E531hP1ttzQRq_OuXITt8Z5BXTXxHNJ0p_h6jYPmvfi6B7w==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:05:48 +0000
X-Pm-External-Id: <CLWD3JfMsFMlyHiwC_RgoqxPiwxaIrtmeiifUCl3DJTcXvTMM7HBkQ4BYllYrDkhR26WI8K5G3gqxGNCjzVIMlAXSeHWNIFgEsPkuuYGG_k=@michelon.ch>
X-Pm-Internal-Id: Ve78Yi7YkjntEfYQ4OfWCzJMKmJBCbw-Ybip7P4E531hP1ttzQRq_OuXITt8Z5BXTXxHNJ0p_h6jYPmvfi6B7w==
To: "dev@michelon.ch" <dev@michelon.ch>
Reply-To: "Test @ Michelon" <test@michelon.ch>
From: "Test @ Michelon" <test@michelon.ch>
Subject: Act 4 - scene 3
X-Pm-Recipient-Encryption: dev%40michelon.ch=pgp-pm
X-Pm-Recipient-Authentication: dev%40michelon.ch=pgp-pm
X-Pm-Scheduled-Sent-Original-Time: Fri, 14 Jun 2024 06:05:34 +0000
Message-Id: <CLWD3JfMsFMlyHiwC_RgoqxPiwxaIrtmeiifUCl3DJTcXvTMM7HBkQ4BYllYrDkhR26WI8K5G3gqxGNCjzVIMlAXSeHWNIFgEsPkuuYGG_k=@michelon.ch>
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:05:48 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;"><span>LODOV=
ICO</span><br><span>I do beseech you, sir, trouble yourself no further.</sp=
an><br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">

            </div>

            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">

            </div>
</div>
)");

    writeFile(dir / "Ve78Yi7YkjntEfYQ4OfWCzJMKmJBCbw-Ybip7P4E531hP1ttzQRq_OuXITt8Z5BXTXxHNJ0p_h6jYPmvfi6B7w==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "Ve78Yi7YkjntEfYQ4OfWCzJMKmJBCbw-Ybip7P4E531hP1ttzQRq_OuXITt8Z5BXTXxHNJ0p_h6jYPmvfi6B7w==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "2",
      "5",
      "7",
      "15"
    ],
    "ExternalID": "CLWD3JfMsFMlyHiwC_RgoqxPiwxaIrtmeiifUCl3DJTcXvTMM7HBkQ4BYllYrDkhR26WI8K5G3gqxGNCjzVIMlAXSeHWNIFgEsPkuuYGG_k=@michelon.ch",
    "Subject": "Act 4 - scene 3",
    "Sender": {
      "Name": "Test @ Michelon",
      "Address": "test@michelon.ch"
    },
    "ToList": [
      {
        "Name": "dev@michelon.ch",
        "Address": "dev@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "Flags": 8206,
    "Time": 1718345148,
    "Size": 714,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 4 - scene 3\r\nTo: dev@michelon.ch \u003cdev@michelon.ch\u003e\r\nFrom: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:05:48 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003cCLWD3JfMsFMlyHiwC_RgoqxPiwxaIrtmeiifUCl3DJTcXvTMM7HBkQ4BYllYrDkhR26WI8K5G3gqxGNCjzVIMlAXSeHWNIFgEsPkuuYGG_k=@michelon.ch\u003e\r\nX-Pm-Scheduled-Sent-Original-Time: Fri, 14 Jun 2024 06:05:34 +0000\r\nX-Pm-Recipient-Authentication: dev%40michelon.ch=pgp-pm\r\nX-Pm-Recipient-Encryption: dev%40michelon.ch=pgp-pm\r\n",
    "WriterType": 0
  }
})");

    writeFile(dir / "ZrVs4o8qNx_xSwO5MtIwNAX__QIqCm3jLVNNJF6jeV3NCgEY-sTzuRaArlKyx0Qx6mpFKp7RB5kqN7CJE_xTrQ==.eml",
              R"(Content-Type: multipart/mixed;
 boundary=a970ebcf334128759e91a825abd6be5d8a9775647e5ff61eec5af3ff305394b5
References: <ZrVs4o8qNx_xSwO5MtIwNAX__QIqCm3jLVNNJF6jeV3NCgEY-sTzuRaArlKyx0Qx6mpFKp7RB5kqN7CJE_xTrQ==@protonmail.internalid>
X-Pm-Date: Fri, 14 Jun 2024 06:10:11 +0000
X-Pm-External-Id: <jr8pUWLzhNpZuAus58i8INLgll5YaFgUZYPKRetqffb2nGLoYOV6ZLkchK9_86mzebhts8zBH6B-Jz1rIz0MhXFniy66igs49MFKgiiC_j8=@michelon.ch>
X-Pm-Internal-Id: ZrVs4o8qNx_xSwO5MtIwNAX__QIqCm3jLVNNJF6jeV3NCgEY-sTzuRaArlKyx0Qx6mpFKp7RB5kqN7CJE_xTrQ==
To: "Test @ Michelon" <test@michelon.ch>
Reply-To: "Dev" <dev@michelon.ch>
From: "Dev" <dev@michelon.ch>
Subject: Mona Lisa.
Delivered-To: test@michelon.ch
Return-Path: <dev@michelon.ch>
X-Original-To: test@michelon.ch
Received: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:10:14 +0000
X-Pm-Spamscore: 0
Message-Id: <jr8pUWLzhNpZuAus58i8INLgll5YaFgUZYPKRetqffb2nGLoYOV6ZLkchK9_86mzebhts8zBH6B-Jz1rIz0MhXFniy66igs49MFKgiiC_j8=@michelon.ch>
X-Attached: MonaLisa.jpg
Mime-Version: 1.0
Date: Fri, 14 Jun 2024 06:10:11 +0000
X-Pm-Origin: internal
X-Pm-Content-Encryption: end-to-end

--a970ebcf334128759e91a825abd6be5d8a9775647e5ff61eec5af3ff305394b5
Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8

<div style=3D"font-family: Arial, sans-serif; font-size: 14px;">Here is Mon=
a Lisa.<br></div>
<div class=3D"protonmail_signature_block protonmail_signature_block-empty" =
style=3D"font-family: Arial, sans-serif; font-size: 14px;">
    <div class=3D"protonmail_signature_block-user protonmail_signature_bloc=
k-empty">
       =20
            </div>
   =20
            <div class=3D"protonmail_signature_block-proton protonmail_sign=
ature_block-empty">
       =20
            </div>
</div>

--a970ebcf334128759e91a825abd6be5d8a9775647e5ff61eec5af3ff305394b5
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename=MonaLisa.jpg
Content-Type: image/jpeg; filename=MonaLisa.jpg; name=MonaLisa.jpg
x-pm-content-encryption: end-to-end

/9j/4AAQSkZJRgABAQEASABIAAD/4QDURXhpZgAATU0AKgAAAAgABwESAAMAAAABAAEAAAEaAAUA
AAABAAAAYgEbAAUAAAABAAAAagEoAAMAAAABAAIAAAExAAIAAAAcAAAAcgEyAAIAAAAUAAAAjodp
AAQAAAABAAAAogAAAAAAAABIAAAAAQAAAEgAAAABQWRvYmUgUGhvdG9zaG9wIENTNSBXaW5kb3dz
ADIwMjQ6MDY6MTQgMDg6MDk6MDUAAAOgAQADAAAAAQABAACgAgADAAAAAQCGAACgAwADAAAAAQDI
AAAAAAAA/+EOamh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8APD94cGFja2V0IGJlZ2luPSLv
u78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4gPHg6eG1wbWV0YSB4bWxuczp4PSJh
ZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgNS41LjAiPiA8cmRmOlJERiB4bWxuczpy
ZGY9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMiPiA8cmRmOkRl
c2NyaXB0aW9uIHJkZjphYm91dD0iIiB4bWxuczpkYz0iaHR0cDovL3B1cmwub3JnL2RjL2VsZW1l
bnRzLzEuMS8iIHhtbG5zOnBob3Rvc2hvcD0iaHR0cDovL25zLmFkb2JlLmNvbS9waG90b3Nob3Av
MS4wLyIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBN
TT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RFdnQ9Imh0dHA6Ly9u
cy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZUV2ZW50IyIgZGM6Zm9ybWF0PSJpbWFn
ZS9qcGVnIiBwaG90b3Nob3A6Q29sb3JNb2RlPSIzIiBwaG90b3Nob3A6SUNDUHJvZmlsZT0ic1JH
QiBJRUM2MTk2Ni0yLjEiIHhtcDpDcmVhdGVEYXRlPSIyMDExLTA2LTA2VDEyOjQwOjA4LTA3OjAw
IiB4bXA6Q3JlYXRvclRvb2w9IkFkb2JlIFBob3Rvc2hvcCBDUzUgV2luZG93cyIgeG1wOk1ldGFk
YXRhRGF0ZT0iMjAyNC0wNi0xNFQwODowOTowNSswMjowMCIgeG1wOk1vZGlmeURhdGU9IjIwMjQt
MDYtMTRUMDg6MDk6MDUrMDI6MDAiIHhtcE1NOkRvY3VtZW50SUQ9InhtcC5kaWQ6RTkzNEU5N0Q4
RDkxRTAxMUJBQzJGRTY4MEY4ODQxQUUiIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6RUEzNEU5
N0Q4RDkxRTAxMUJBQzJGRTY4MEY4ODQxQUUiIHhtcE1NOk9yaWdpbmFsRG9jdW1lbnRJRD0ieG1w
LmRpZDpFOTM0RTk3RDhEOTFFMDExQkFDMkZFNjgwRjg4NDFBRSI+IDx4bXBNTTpIaXN0b3J5PiA8
cmRmOlNlcT4gPHJkZjpsaSB4bXBNTTphY3Rpb249ImNyZWF0ZWQiIHhtcE1NOmluc3RhbmNlSUQ9
InhtcC5paWQ6RTkzNEU5N0Q4RDkxRTAxMUJBQzJGRTY4MEY4ODQxQUUiIHhtcE1NOnNvZnR3YXJl
QWdlbnQ9IkFkb2JlIFBob3Rvc2hvcCBDUzUgV2luZG93cyIgeG1wTU06d2hlbj0iMjAxMS0wNi0w
NlQxMjo0MDowOC0wNzowMCIvPiA8cmRmOmxpIHhtcE1NOmFjdGlvbj0iY29udmVydGVkIiB4bXBN
TTpwYXJhbWV0ZXJzPSJmcm9tIGltYWdlL3RpZmYgdG8gaW1hZ2UvanBlZyIvPiA8cmRmOmxpIHht
cE1NOmFjdGlvbj0ic2F2ZWQiIHhtcE1NOmNoYW5nZWQ9Ii8iIHhtcE1NOmluc3RhbmNlSUQ9Inht
cC5paWQ6RUEzNEU5N0Q4RDkxRTAxMUJBQzJGRTY4MEY4ODQxQUUiIHhtcE1NOnNvZnR3YXJlQWdl
bnQ9IkFkb2JlIFBob3Rvc2hvcCBDUzUgV2luZG93cyIgeG1wTU06d2hlbj0iMjAxMS0wNi0wN1Qy
MjowOTozMC0wNzowMCIvPiA8cmRmOmxpIHN0RXZ0OmFjdGlvbj0icHJvZHVjZWQiIHN0RXZ0OnNv
ZnR3YXJlQWdlbnQ9IkFmZmluaXR5IFBob3RvIDIgMi41LjIiIHN0RXZ0OndoZW49IjIwMjQtMDYt
MTRUMDg6MDk6MDUrMDI6MDAiLz4gPC9yZGY6U2VxPiA8L3htcE1NOkhpc3Rvcnk+IDwvcmRmOkRl
c2NyaXB0aW9uPiA8L3JkZjpSREY+IDwveDp4bXBtZXRhPiAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAg
ICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIDw/eHBhY2tldCBlbmQ9InciPz7/7QAsUGhv
dG9zaG9wIDMuMAA4QklNBCUAAAAAABDUHYzZjwCyBOmACZjs+EJ+/+ICZElDQ19QUk9GSUxFAAEB
AAACVGxjbXMEMAAAbW50clJHQiBYWVogB+gABgANAAYAEQAeYWNzcEFQUEwAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAPbWAAEAAAAA0y1sY21zAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAALZGVzYwAAAQgAAAA+Y3BydAAAAUgAAABMd3RwdAAAAZQAAAAUY2hh
ZAAAAagAAAAsclhZWgAAAdQAAAAUYlhZWgAAAegAAAAUZ1hZWgAAAfwAAAAUclRSQwAAAhAAAAAg
Z1RSQwAAAhAAAAAgYlRSQwAAAhAAAAAgY2hybQAAAjAAAAAkbWx1YwAAAAAAAAABAAAADGVuVVMA
AAAiAAAAHABzAFIARwBCACAASQBFAEMANgAxADkANgA2AC0AMgAuADEAAG1sdWMAAAAAAAAAAQAA
AAxlblVTAAAAMAAAABwATgBvACAAYwBvAHAAeQByAGkAZwBoAHQALAAgAHUAcwBlACAAZgByAGUA
ZQBsAHlYWVogAAAAAAAA9tYAAQAAAADTLXNmMzIAAAAAAAEMQgAABd7///MlAAAHkwAA/ZD///uh
///9ogAAA9wAAMBuWFlaIAAAAAAAAG+gAAA49QAAA5BYWVogAAAAAAAAJJ8AAA+EAAC2w1hZWiAA
AAAAAABilwAAt4cAABjZcGFyYQAAAAAAAwAAAAJmZgAA8qcAAA1ZAAAT0AAACltjaHJtAAAAAAAD
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
--a970ebcf334128759e91a825abd6be5d8a9775647e5ff61eec5af3ff305394b5--
)");

    writeFile(dir / "ZrVs4o8qNx_xSwO5MtIwNAX__QIqCm3jLVNNJF6jeV3NCgEY-sTzuRaArlKyx0Qx6mpFKp7RB5kqN7CJE_xTrQ==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "ZrVs4o8qNx_xSwO5MtIwNAX__QIqCm3jLVNNJF6jeV3NCgEY-sTzuRaArlKyx0Qx6mpFKp7RB5kqN7CJE_xTrQ==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "0",
      "5",
      "15"
    ],
    "ExternalID": "jr8pUWLzhNpZuAus58i8INLgll5YaFgUZYPKRetqffb2nGLoYOV6ZLkchK9_86mzebhts8zBH6B-Jz1rIz0MhXFniy66igs49MFKgiiC_j8=@michelon.ch",
    "Subject": "Mona Lisa.",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718345411,
    "Size": 9936,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 1,
    "Attachments": [
      {
        "ID": "dMR8YTLiutaxLrK0J2etXJ9wRavvvg11-8i97jLfS1w15aylLeZ31f6dqdQjVskh6bStEgyqYJ4fEqbcMV_p5g==",
        "Name": "MonaLisa.jpg",
        "Size": 9273,
        "MIMEType": "image/jpeg",
        "Disposition": "attachment",
        "Headers": {
          "content-disposition": "attachment",
          "x-pm-content-encryption": "end-to-end"
        },
        "KeyPackets": "wV4DR2zSGHm0e5QSAQdAARwzSmy1qvq+VjiYRrlFEMhbY6/nYvywDONKOob+Xy0wp1gSfQcXnEOOYmdUBMpq2k5n3ChblDipzRfHSJN17BO7+HT+Y2kz2PLaNR3h3QTi",
        "Signature": "-----BEGIN PGP SIGNATURE-----\nVersion: ProtonMail\n\nwnUEABYKACcFgmZr3qYJkLVY6yXWhsidFiEEo9tis1m9Sj3juKf9tVjrJdaG\nyJ0AALyAAP4kGx3rPkgsJUlWkKArIw147Sf/AuAOfxL5eRQCSVI8oQEA+UIf\nsIylJ98IPWgiBGrlKebn+osWF1Q2jPTuBmlfQQA=\n=TPO8\n-----END PGP SIGNATURE-----\n"
      }
    ],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Mona Lisa.\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:10:11 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nX-Attached: MonaLisa.jpg\r\nMessage-Id: \u003cjr8pUWLzhNpZuAus58i8INLgll5YaFgUZYPKRetqffb2nGLoYOV6ZLkchK9_86mzebhts8zBH6B-Jz1rIz0MhXFniy66igs49MFKgiiC_j8=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:10:14 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");
}

//****************************************************************************************************************************************************
/// \param[in] dir The folder path. The folder must exist and contain an existing backup.
//****************************************************************************************************************************************************
void addSkippedAndFailingMessages(std::filesystem::path const& dir) {
    writeFile(dir / "jBnHMoZ9R1sK27iVRXkc-o81MUf3s4yamEcouHl5RiEmSx19bgcw7hm884h3LMsamTIuuBKn5dBHgjBGiaDVBw==.eml", R"()");
    writeFile(dir / "jBnHMoZ9R1sK27iVRXkc-o81MUf3s4yamEcouHl5RiEmSx19bgcw7hm884h3LMsamTIuuBKn5dBHgjBGiaDVBw==.metadata.json", R"({
  "Version": 1,
  "Payload": {
    "ID": "non-exiting ID",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "0",
      "5",
      "15"
    ],
    "ExternalID": "jr8pUWLzhNpZuAus58i8INLgll5YaFgUZYPKRetqffb2nGLoYOV6ZLkchK9_86mzebhts8zBH6B-Jz1rIz0MhXFniy66igs49MFKgiiC_j8=@michelon.ch",
    "Subject": "Mona Lisa.",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718345411,
    "Size": 9936,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Mona Lisa.\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:10:11 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nX-Attached: MonaLisa.jpg\r\nMessage-Id: \u003cjr8pUWLzhNpZuAus58i8INLgll5YaFgUZYPKRetqffb2nGLoYOV6ZLkchK9_86mzebhts8zBH6B-Jz1rIz0MhXFniy66igs49MFKgiiC_j8=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:10:14 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");
    writeFile(dir / "2yVSh86fKc_VnjuFUZTxTFBAJBqoOhTnpusBR29e38a6sKJb4FYe1WDZZaY-qRxt9-NJoyK3FlDb41iKKSdPGg==.eml",
              std::string(1 << 20, '0')); // this will trigger "header exceeds maximum size".

    writeFile(dir / "2yVSh86fKc_VnjuFUZTxTFBAJBqoOhTnpusBR29e38a6sKJb4FYe1WDZZaY-qRxt9-NJoyK3FlDb41iKKSdPGg==.metadata.json",
              R"({
  "Version": 1,
  "Payload": {
    "ID": "2yVSh86fKc_VnjuFUZTxTFBAJBqoOhTnpusBR29e38a6sKJb4FYe1WDZZaY-qRxt9-NJoyK3FlDb41iKKSdPGg==",
    "AddressID": "o1NktWYlte9s3miWOwimQSV1nVsJX36FlBRfX25E48SPeHTSO2YrWn3K5nd9cR-1jdmxkxExZ3DlaiNnk7XA2g==",
    "LabelIDs": [
      "5",
      "15",
      "UD2CCpFokJimuXRGM91VglWoj4eOLCevCs0JIfVDcmlLmlqUxVcqWW_6UMyz7D17XBbcqK4SDxUBO66KY8LlCw=="
    ],
    "ExternalID": "RA_0QgwMkjDSg5TTCv93_uGZkmWBFwys-fY9Rn_fLCv_3GaRNq2czzyc4F8SqmSQqL_FxHk_8iku3Heb5h2g0l4L0xMk6I9hJNDYf6KkqIM=@michelon.ch",
    "Subject": "Act 3 - Scene 1",
    "Sender": {
      "Name": "Dev",
      "Address": "dev@michelon.ch"
    },
    "ToList": [
      {
        "Name": "Test @ Michelon",
        "Address": "test@michelon.ch"
      }
    ],
    "CCList": [],
    "BCCList": [],
    "ReplyTos": [
      {
        "Name": "Dev",
        "Address": "dev@michelon.ch"
      }
    ],
    "Flags": 9229,
    "Time": 1718344974,
    "Size": 380,
    "Unread": 0,
    "IsReplied": 0,
    "IsRepliedAll": 0,
    "IsForwarded": 0,
    "NumAttachments": 0,
    "Attachments": [],
    "MIMEType": "text/html",
    "Headers": "X-Pm-Content-Encryption: end-to-end\r\nX-Pm-Origin: internal\r\nSubject: Act 3 - Scene 1\r\nTo: Test @ Michelon \u003ctest@michelon.ch\u003e\r\nFrom: Dev \u003cdev@michelon.ch\u003e\r\nDate: Fri, 14 Jun 2024 06:02:54 +0000\r\nMime-Version: 1.0\r\nContent-Type: text/html\r\nMessage-Id: \u003cRA_0QgwMkjDSg5TTCv93_uGZkmWBFwys-fY9Rn_fLCv_3GaRNq2czzyc4F8SqmSQqL_FxHk_8iku3Heb5h2g0l4L0xMk6I9hJNDYf6KkqIM=@michelon.ch\u003e\r\nX-Pm-Spamscore: 0\r\nReceived: from mail.protonmail.ch by mail.protonmail.ch; Fri, 14 Jun 2024 06:02:59 +0000\r\nX-Original-To: test@michelon.ch\r\nReturn-Path: \u003cdev@michelon.ch\u003e\r\nDelivered-To: test@michelon.ch\r\n",
    "WriterType": 0
  }
})");
}
