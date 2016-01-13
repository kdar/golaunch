' GetRealTarget.vbs
' This version needs to be run under wscript engine rather than cscript

' Pass the full path to an MSI "Advertised Shortcut" lnk file (including the extension) as a parameter
' e.g. assuming that we have a default install of Office 2003 for All Users:
' GetRealTarget "C:\Documents and Settings\All Users\Start Menu\Programs\Microsoft Office\Microsoft Office Excel 2003.lnk" 
' Displays fully resolved target for the MSI shortcut

Option Explicit
Dim MSITarget

On Error Resume Next ' just some simple error handling for purposes of this example
If wscript.arguments.count = 1 Then ' did actually pass an MSI advertised shortcut? Or, at least, a parameter that could be such a thing?
   With CreateObject("WindowsInstaller.Installer")
      Set MSITarget = .ShortcutTarget(wscript.arguments(0))
      If Err = 0 then
         MsgBox .ComponentPath(MSITarget.StringData(1), MSITarget.StringData(3))
      Else 
         MsgBox wscript.arguments(0) & vbcrlf & "is not a legitimate MSI shortcut file or could not be found"
      End If
   End With
End If
On Error Goto 0