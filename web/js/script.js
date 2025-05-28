document.addEventListener("DOMContentLoaded", () => {
  const cart = [];
  const cartIcon = document.getElementById("cart-icon");
  const cartCount = document.getElementById("cart-count");
  const cartModal = document.getElementById("cart-modal");
  const cartContent = document.getElementById("cart-content");
  const checkoutButton = document.getElementById("checkout-button");
  const closeCartButton = document.getElementById("close-cart");

  // корзина в localstorage
  const saveCartToStorage = () => {
    localStorage.setItem("cart", JSON.stringify(cart));
  };

  // загрузка корзины
  const loadCartFromStorage = () => {
    const storedCart = localStorage.getItem("cart");
    if (storedCart) {
      const parsedCart = JSON.parse(storedCart);
      cart.push(...parsedCart);
    }
  };

  // Функция обновления иконки корзины
  const updateCartIcon = () => {
    const totalItems = cart.reduce((sum, item) => sum + item.quantity, 0);
    cartCount.textContent = totalItems;
    cartIcon.style.display = "flex";
  };

  // Функция обновления модального окна корзины
  const updateCartModal = () => {
    if (cart.length === 0) {
      cartContent.innerHTML = "<p>Ваша корзина пуста.</p>";
      checkoutButton.style.display = "none";
    } else {
      checkoutButton.style.display = "block";
      const total = cart.reduce(
        (sum, item) => sum + item.price * item.quantity,
        0,
      );
      const cartItemsHTML = cart
        .map(
          (item) => `
                <div class="cart-item">
                <span>${item.quantity} шт. ${item.name} </span>
                <span>${item.price} руб</span>
                    <button class="remove-item" data-name="${item.name}">Удалить</button>
                </div>`,
        )
        .join("");

      cartContent.innerHTML = `
                <div class="cart-items">${cartItemsHTML}</div>
                <div class="cart-total"><strong>Итого:</strong> ${total} руб.</div>`;
    }
  };

  //добавление товара в корзину
  const addToCart = (name, price) => {
    const existingItem = cart.find((item) => item.name === name);
    if (existingItem) {
      existingItem.quantity += 1;
    } else {
      cart.push({ name, price: parseFloat(price), quantity: 1 });
    }
    saveCartToStorage();
    updateCartIcon();
    updateCartModal();
  };

  // Функция удаления товара из корзины
  const removeFromCart = (name) => {
    const itemIndex = cart.findIndex((item) => item.name === name);
    if (itemIndex > -1) {
      cart.splice(itemIndex, 1);
      saveCartToStorage(); // Сохраняем корзину
      updateCartIcon();
      updateCartModal();
    }
  };

  // Обработчик кликов для добавления в корзину
  document.body.addEventListener("click", (event) => {
    if (event.target.classList.contains("add-to-cart")) {
      event.preventDefault();
      const name = event.target.dataset.name;
      const price = event.target.dataset.price;
      addToCart(name, price);
    }
  });

  // Обработчик кликов для удаления из корзины
  cartContent.addEventListener("click", (event) => {
    if (event.target.classList.contains("remove-item")) {
      const name = event.target.dataset.name;
      removeFromCart(name);
    }
  });

  // Обработчик открытия и закрытия корзины
  cartIcon.addEventListener("click", () => {
    cartModal.style.display =
      cartModal.style.display === "none" ? "block" : "none";
  });

  closeCartButton.addEventListener("click", () => {
    cartModal.style.display = "none";
  });

  // Загружаем корзину при старте
  loadCartFromStorage();
  updateCartIcon();
  updateCartModal();
});
