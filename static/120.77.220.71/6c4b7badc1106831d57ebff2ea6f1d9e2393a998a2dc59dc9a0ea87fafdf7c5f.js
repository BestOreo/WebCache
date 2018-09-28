var board = document.getElementById("content")

function showname(str) {
    console.log(str)
    var p = document.createElement("span");
    if (str == "I") {
        p.className = "text-default"
    } else if (str == "Like") {
        p.style.color = "#337ab7"
    } else if (str == "University") {
        p.style.color = "#5cb85c"
    } else if (str == "of") {
        p.style.color = "#5bc0de"
    } else if (str == "British") {
        p.style.color = "#f0ad4e"
    } else if (str == "Columbia") {
        p.style.color = "#d9534f"
    }
    p.textContent = str + "\t"
    board.appendChild(p)
}