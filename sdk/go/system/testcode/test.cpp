#include <iostream>
#include <windows.h>
#include <Msi.h>

using namespace std;

int main(int argc, char *argv[]) {
  CoInitializeEx(0, COINIT_APARTMENTTHREADED | COINIT_SPEED_OVER_MEMORY);
  //CoInitialize(0);
  wchar_t szProductCode[MAX_GUID_CHARS + 1] = {0};
	wchar_t szFeatureId[MAX_FEATURE_CHARS + 1] = {0};
	wchar_t szComponentCode[MAX_GUID_CHARS + 1] = {0};
  wchar_t szPath[MAX_PATH + 1] = {0};
  DWORD dwLen = MAX_PATH;

  cout << MsiGetShortcutTargetW(L"..\\testdata\\Minesweeper.lnk", szProductCode, szFeatureId, szComponentCode) << endl;
  cout << MsiGetComponentPathW(szProductCode, szComponentCode, szPath, &dwLen) << endl;
  wcout << szPath << endl;
}
