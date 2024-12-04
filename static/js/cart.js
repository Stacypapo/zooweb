document.addEventListener("DOMContentLoaded", function () {
    const removeButtons = document.querySelectorAll(".remove-btn");

    removeButtons.forEach((button) => {
        button.addEventListener("click", function () {
            const cartItemId = this.dataset.id; // ID товара из data-id

            fetch("/user/cart/remove", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ cart_item_id: cartItemId }), // Передаем ID товара
            })
                .then((response) => {
                    if (!response.ok) {
                        throw new Error("Ошибка при удалении товара");
                    }
                    return response.json();
                })
                .then((data) => {
                    if (data.message === "Товар удален из корзины") {
                        // Удаляем строку товара из таблицы
                        const row = document.getElementById(`cart-item-${cartItemId}`);
                        if (row) {
                            row.remove();
                        }
                    } else {
                        alert("Не удалось удалить товар из корзины.");
                    }
                })
                .catch((error) => {
                    console.error("Ошибка:", error);
                    alert("Произошла ошибка. Попробуйте позже.");
                });
        });
    });
});
