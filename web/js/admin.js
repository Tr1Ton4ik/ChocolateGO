const tabs = document.querySelectorAll("nav a");
const sections = document.querySelectorAll("main > section");

tabs.forEach((tab, index) => {
  tab.addEventListener("click", (e) => {
    e.preventDefault();
    sections.forEach((section) => (section.style.display = "none"));
    sections[index].style.display = "block";
  });
});

function getCookie(name) {
  let cookie = document.cookie
    .split("; ")
    .find((row) => row.startsWith(name + "="));
  return cookie ? cookie.split("=")[1] : null;
}

if (getCookie("isAuthenticated") !== "true") {
  window.location.href = "/login";
}

// Обработчик клика для кнопки выхода
document.getElementById("logout-btn").addEventListener("click", (event) => {
  event.preventDefault();
  document.cookie = "isAuthenticated=false; path=/";
  window.location.href = "/login";
});
document.getElementById("admin-form").addEventListener("submit", (event) => {
  event.preventDefault();

  const password = document.getElementById("password").value;
  const confirmPassword = document.getElementById("confirm-password").value;

  if (password !== confirmPassword) {
    alert("Пароли не совпадают. Пожалуйста, проверьте еще раз.");
    return;
  }

  // Если пароли совпадают, отправляем форму
  event.target.submit();
});
async function searchProduct() {
  const query = document.getElementById("search-input").value;
  const response = await fetch(
    `/search_product?name=${encodeURIComponent(query)}`,
  );
  const results = await response.json();
  displayResults(results);
}

// Функция для отображения результатов поиска
function displayResults(results) {
  const resultsContainer = document.getElementById("search-results");
  resultsContainer.innerHTML = ""; // очищаем предыдущие результаты
  if (results.length > 0) {
    results.forEach((product) => {
      const div = document.createElement("div");
      div.textContent = product.name + `(${product.category})`; // Предполагаем, что у товара есть свойство 'name'
      resultsContainer.appendChild(div);
    });
  } else {
    resultsContainer.textContent = "Нет результатов";
  }
}
document.addEventListener("DOMContentLoaded", function () {
  const productForm = document.getElementById("product-form");

  if (productForm) {
    productForm.addEventListener("submit", function (event) {
      const priceInput = document.getElementById("price");
      const priceValue = parseFloat(priceInput.value);

      if (priceValue <= 0 || isNaN(priceValue)) {
        event.preventDefault();
        alert("Цена должна быть больше 0. Пожалуйста, исправьте значение.");
        priceInput.focus();
      }
    });
  }
});

document
  .getElementById("delete-category-form")
  .addEventListener("submit", function (event) {
    event.preventDefault();

    const categoryName = document.getElementById("category-to-delete").value;
    const csrfToken = document.querySelector(
      'input[name="gorilla.csrf.Token"]',
    ).value;

    if (categoryName) {
      fetch("/delete_category", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": csrfToken,
        },
        body: JSON.stringify({ name: categoryName }),
      }).then((response) => {
        if (response.ok) {
          window.location.reload();
        }
      });
    }
  });

document
  .getElementById("delete-form")
  .addEventListener("submit", function (event) {
    event.preventDefault();

    const searchName = document.getElementById("search-input").value;
    const productName = document.getElementById("delete-product-name").value;
    const categoryName = document.getElementById(
      "delete-product-category",
    ).value;
    const csrfToken = document.querySelector(
      'input[name="gorilla.csrf.Token"]',
    ).value;

    if (categoryName && searchName == productName) {
      fetch("/delete_product", {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": csrfToken,
        },
        body: JSON.stringify({
          category: categoryName,
          product: productName,
        }),
      }).then((response) => {
        if (response.ok) {
          window.location.reload();
        }
      });
    }
  });
