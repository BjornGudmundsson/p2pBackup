console.log(5);

document.getElementById("backupButton").addEventListener("click", ()=> {
    let fn = document.getElementById("backupFile").value;
    console.log(fn);
})

document.getElementById("removeButton").addEventListener("click", () => {
    fn = document.getElementById("removeFile").value;
    console.log(fn);
})