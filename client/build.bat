@echo off
echo Building CLI client...
cd cli
SET CGO_ENABLED=0
"C:\Program Files\Go\bin\go.exe" build -o ..\SushiSyncCLI.exe
cd ..

echo Building GUI client (may fail if CGO dependencies are missing)...
SET CGO_ENABLED=0
"C:\Program Files\Go\bin\go.exe" build -tags nocgo -o SushiSyncGUI.exe

echo Done! Check for SushiSyncCLI.exe and SushiSyncGUI.exe
