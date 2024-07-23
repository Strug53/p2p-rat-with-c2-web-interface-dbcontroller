const term = new Terminal();
const terminalContainer = document.getElementById('terminal');
term.open(terminalContainer);
term.write('Welcome to Xterm.js Demo\n \r');
PrintNewRaw()

let command = '';

term.onData( e => {
    if (e === '\r') { // Enter key
        term.write("\n \r")
        console.log(command)
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
async function SendCommand(){

    let data = {
        id: "1",
        cmd: "ls"
    }
    console.log(data)
    const response = await fetch("/sendCmd", {
        method: "POST",
        headers: {'Content-Type': 'application/json'}, 
        body: JSON.stringify(data)
    })
    const body = await response.text()
    term.write(body)
    
};

async function handleCommand(input) {
    const args = input.split(' ');
    const command = args[0];
    const params = args.slice(1).join(' ');

    switch (command) {
        case 'clear':
            term.clear();
            break;
        case 'cls':
            term.clear()
            break;
        default:
            await SendCommand();
            break;
    }
    PrintNewRaw()
}

// clear command
//cls
