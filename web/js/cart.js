document.addEventListener("DOMContentLoaded", function () {
  // Функция для пересчёта общей суммы
  function updateTotal() {
    let total = 0;
    const cartItems = document.querySelectorAll(".cart-item");
    cartItems.forEach((item) => {
      const priceElement = item.querySelector(".price");
      let price = 0;

      if (priceElement.dataset.price) {
        price = parseFloat(priceElement.dataset.price);
      } else {
        price = parseFloat(priceElement.textContent.replace(/[^\d.]/g, ""));
      }
      const quantity = parseInt(
        item.querySelector('input[type="number"]').value,
      );
      total += price * quantity;
    });
    document.querySelector(".total-price").textContent = total + " ₽";
  }

  // Обработчик увеличения количества
  document.querySelectorAll(".increase").forEach((button) => {
    button.addEventListener("click", function () {
      const input = this.parentElement.querySelector('input[type="number"]');
      let quantity = parseInt(input.value);
      input.value = quantity + 1;
      updateTotal();
    });
  });

  // Обработчик уменьшения количества
  document.querySelectorAll(".decrease").forEach((button) => {
    button.addEventListener("click", function () {
      const input = this.parentElement.querySelector('input[type="number"]');
      let quantity = parseInt(input.value);
      if (quantity > 1) {
        input.value = quantity - 1;
        updateTotal();
      }
    });
  });

  // Обработчик удаления товара из корзины
  document.querySelectorAll(".remove-item").forEach((button) => {
    button.addEventListener("click", function () {
      this.closest(".cart-item").remove();
      updateTotal();
    });
  });

  // При изменении количества через ручной ввод также обновляем итог
  document.querySelectorAll('input[type="number"]').forEach((input) => {
    input.addEventListener("change", updateTotal);
  });
});
