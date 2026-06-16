document.addEventListener("DOMContentLoaded", function () {
    var avatarWrap = document.getElementById("avatar-wrap");
    var avatarInput = document.getElementById("avatar-input");
    var avatarForm = document.getElementById("avatar-form");
    var deleteBtn = document.getElementById("delete-account-btn");

    if (avatarWrap && avatarInput) {
        avatarWrap.addEventListener("click", function () {
            avatarInput.click();
        });

        avatarInput.addEventListener("change", function () {
            if (avatarInput.files.length === 0) return;

            var formData = new FormData(avatarForm);

            fetch("/api/profile/avatar", {
                method: "POST",
                body: formData
            })
                .then(function (r) { return r.json(); })
                .then(function (data) {
                    if (data.success) {
                        window.location.reload();
                    } else {
                        alert(data.message || "Erreur");
                    }
                })
                .catch(function (err) {
                    alert("Erreur lors de l'upload");
                });
        });
    }

    if (deleteBtn) {
        deleteBtn.addEventListener("click", function () {
            if (!confirm("Es-tu sûr de vouloir supprimer ton compte ? Cette action est irréversible.")) {
                return;
            }

            fetch("/api/profile/delete", {
                method: "POST"
            })
                .then(function (r) { return r.json(); })
                .then(function (data) {
                    if (data.success) {
                        window.location.href = "/";
                    } else {
                        alert(data.message || "Erreur");
                    }
                })
                .catch(function () {
                    alert("Erreur serveur");
                });
        });
    }
});
