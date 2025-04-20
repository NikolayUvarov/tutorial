@echo off
setlocal enabledelayedexpansion


set "TARGET_DIR=test-set"
set "ZIP_PATTERN=*.zip"

if exist "%TARGET_DIR%\" (
    echo Exitsts "%TARGET_DIR%", proceed with tests
    goto :MAINTEST
)

echo Folder "%TARGET_DIR%" not found, creating...
mkdir "%TARGET_DIR%"

for %%F in (%ZIP_PATTERN%) do (
    set "ZIP_FILE=%%F"
    goto :UNZIP
)

echo No zip-files found, exitting. You can create %TARGET_DIR% manually if you have data for tests.
goto :EOF

:UNZIP
echo Found zip-file: "%ZIP_FILE%"
REM Распаковка PowerShell или tar
powershell -Command "Expand-Archive -LiteralPath '%ZIP_FILE%' -DestinationPath '%TARGET_DIR%'"
REM tar -xf "%ZIP_FILE%" -C "%TARGET_DIR%"

echo Extract done.

:MAINTEST

set P=

echo Loop through the items in the current directory - looking for first directory, assuming there is only one with test set 
for /d %%i in (*) do (
    echo Curr: %%i
    rem If the variable P is empty, set it to the first directory name
    if "!P!"=="" set P=%%i
)

echo First directory: %P%

go build -o main.exe . 
if %errorlevel% NEQ 0 (
    echo Build failed! Exiting.
    echo ======================
    exit /b 1
)

for %%f in (%P%\*.*) do (
    if "%%~xf" NEQ ".a" if "%%~xf" NEQ ".out" (
        REM Run the Go program, redirect input from test file
        main.exe < "%P%\%%~nf" > "%P%\%%~nf.out"

        echo -----------------------------------------

        REM Compare output with the expected result
        fc "%P%\%%~nf.out" "%P%\%%~nf.a" > nul
        if errorlevel 1 (
            echo *** Test case %%~nf: Fail ***
            echo Expected:
            type "%P%\%%~nf.a"
            echo.
            echo Got:
            type "%P%\%%~nf.out"
            timeout 1
	        exit
        ) else (
            echo Test case %%~nf: Pass
        )

        REM Optionally remove the output file
        REM del "%P%\%%~nf.out"
    )
)

endlocal

echo ALL DONE
timeout /t 5

