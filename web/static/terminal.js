const term = new Terminal();
const terminalContainer = document.getElementById('terminal');
term.open(terminalContainer);
term.write('Welcome to Xterm.js Demo\n \r');
PrintNewRaw()

let command = '';

term.onData(e => {
    if (e === '\r') { // Enter key
        term.write("\n \r")
        handleCommand(command.trim());
        command = ''; // Reset command buffer
    } else if (e === '\x7F') { // Backspace key
        if (command.length > 0) {
            term.write('\b \b'); // Move cursor back, erase character, move cursor back again
            command = command.slice(0, -1); // Remove last character from command buffer
        }
    } else {
        term.write(e); // Echo the typed character
        command += e; // Add typed character to command buffer
    }
});
function PrintNewRaw(){
    term.write("\n\r")
    term.write("# ")

}
function handleCommand(input) {
    const args = input.split(' ');
    const command = args[0];
    const params = args.slice(1).join(' ');

    switch (command) {
        case 'echo':
            term.write(params + '\n \r');
            break;
        case 'cat':
            // Implement 'cat' command logic here
            break;
        default:
            term.write(`Command not found: ${command}\n \r`);
            break;
    }
    PrintNewRaw()
}

// clear command
//cls
