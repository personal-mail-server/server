const loginForm = document.getElementById("login-form");
const loginView = document.getElementById("login-view");
const authView = document.getElementById("auth-view");
const loginIdInput = document.getElementById("loginId");
const passwordInput = document.getElementById("password");
const submitButton = document.getElementById("submit-button");
const logoutButton = document.getElementById("logout-button");
const userMenuLabel = document.getElementById("user-menu-label");
const sessionLoginId = document.getElementById("session-login-id");
const sessionTokenType = document.getElementById("session-token-type");
const sessionAccessExpiry = document.getElementById("session-access-expiry");
const sessionRefreshExpiry = document.getElementById("session-refresh-expiry");
const statusPanel = document.getElementById("status-panel");
const statusMessage = document.getElementById("status-message");
const statusPayload = document.getElementById("status-payload");

const refreshAddressesButton = document.getElementById("refresh-addresses-button");
const addressListMeta = document.getElementById("address-list-meta");
const addressListEmpty = document.getElementById("address-list-empty");
const addressList = document.getElementById("address-list");
const addressDetailEmpty = document.getElementById("address-detail-empty");
const addressDetailCard = document.getElementById("address-detail-card");
const detailEmail = document.getElementById("detail-email");
const detailBadge = document.getElementById("detail-badge");
const detailId = document.getElementById("detail-id");
const detailOwner = document.getElementById("detail-owner");
const detailCreatedAt = document.getElementById("detail-created-at");
const detailUpdatedAt = document.getElementById("detail-updated-at");
const detailDeletedAt = document.getElementById("detail-deleted-at");
const createAddressForm = document.getElementById("create-address-form");
const createEmailInput = document.getElementById("create-email");
const generateCandidateButton = document.getElementById("generate-candidate-button");
const createAddressButton = document.getElementById("create-address-button");
const candidateHelper = document.getElementById("candidate-helper");

const loginIdPattern = /^[a-z0-9-]{4,32}$/;
const passwordPattern = /^(?=.*[A-Za-z])(?=.*[0-9])[^\s]{8,64}$/;
const simpleEmailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

let authState = null;

function setStatus(kind, message, payload) {
  statusPanel.className = `status-panel ${kind}`;
  statusMessage.textContent = message;
  statusPayload.textContent = JSON.stringify(payload, null, 2);
}

function resetAddressViews() {
  addressList.innerHTML = "";
  addressListMeta.textContent = "0개의 주소";
  addressListEmpty.classList.remove("hidden");
  addressDetailEmpty.classList.remove("hidden");
  addressDetailCard.classList.add("hidden");
  detailEmail.textContent = "-";
  detailBadge.textContent = "ACTIVE";
  detailId.textContent = "-";
  detailOwner.textContent = "-";
  detailCreatedAt.textContent = "-";
  detailUpdatedAt.textContent = "-";
  detailDeletedAt.textContent = "-";
  candidateHelper.textContent = "후보값을 받아 입력창을 채우거나, 원하는 주소를 직접 입력하세요.";
}

function renderLoggedOut() {
  loginView.classList.remove("hidden");
  authView.classList.add("hidden");
  createEmailInput.value = "";
  resetAddressViews();
}

function formatDateTime(value) {
  if (!value) {
    return "-";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString("ko-KR", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
}

function renderAddressDetail(address) {
  if (!address) {
    addressDetailEmpty.classList.remove("hidden");
    addressDetailCard.classList.add("hidden");
    return;
  }

  addressDetailEmpty.classList.add("hidden");
  addressDetailCard.classList.remove("hidden");
  detailEmail.textContent = address.email;
  detailBadge.textContent = address.deletedAt ? "DELETED" : "ACTIVE";
  detailId.textContent = String(address.id);
  detailOwner.textContent = String(address.ownerUserId);
  detailCreatedAt.textContent = formatDateTime(address.createdAt);
  detailUpdatedAt.textContent = formatDateTime(address.updatedAt);
  detailDeletedAt.textContent = formatDateTime(address.deletedAt);
}

function getSelectedAddress() {
  if (!authState || !Array.isArray(authState.addresses)) {
    return null;
  }

  return authState.addresses.find((address) => address.id === authState.selectedAddressId) || null;
}

function renderAddressList() {
  const addresses = authState && Array.isArray(authState.addresses) ? authState.addresses : [];

  addressList.innerHTML = "";
  addressListMeta.textContent = `${addresses.length}개의 주소`;

  if (addresses.length === 0) {
    addressListEmpty.classList.remove("hidden");
    renderAddressDetail(null);
    return;
  }

  addressListEmpty.classList.add("hidden");

  addresses.forEach((address) => {
    const item = document.createElement("button");
    item.type = "button";
    item.className = "address-item";
    if (address.id === authState.selectedAddressId) {
      item.classList.add("is-active");
    }
    item.dataset.addressId = String(address.id);
    item.innerHTML = `
      <div class="address-item-top">
        <span class="address-item-email">${address.email}</span>
        <span class="detail-badge">${address.deletedAt ? "DELETED" : "ACTIVE"}</span>
      </div>
      <div class="address-item-meta">
        생성: ${formatDateTime(address.createdAt)}
      </div>
    `;
    addressList.appendChild(item);
  });

  renderAddressDetail(getSelectedAddress());
}

function renderAuthenticated() {
  if (!authState) {
    renderLoggedOut();
    return;
  }

  loginView.classList.add("hidden");
  authView.classList.remove("hidden");
  userMenuLabel.textContent = authState.loginId;
  sessionLoginId.textContent = authState.loginId;
  sessionTokenType.textContent = authState.tokenType;
  sessionAccessExpiry.textContent = String(authState.accessTokenExpiresIn);
  sessionRefreshExpiry.textContent = String(authState.refreshTokenExpiresIn);
  renderAddressList();
}

function clearAuthState() {
  authState = null;
}

async function readJsonIfPossible(response) {
  const contentType = response.headers.get("Content-Type") || "";
  if (!contentType.includes("application/json")) {
    return {};
  }
  return response.json();
}

function validateInputs(loginId, password) {
  if (!loginIdPattern.test(loginId)) {
    return {
      code: "INVALID_LOGIN_ID_FORMAT",
      message: "로그인 ID는 4~32자 영문 소문자, 숫자, 하이픈만 사용할 수 있습니다.",
    };
  }

  if (!passwordPattern.test(password)) {
    return {
      code: "INVALID_PASSWORD_FORMAT",
      message: "비밀번호는 8~64자이며 영문자와 숫자를 각각 1개 이상 포함해야 합니다.",
    };
  }

  return null;
}

function validateCreateEmail(email) {
  if (!email) {
    return {
      code: "MISSING_REQUIRED_FIELD",
      message: "테스트용 메일 주소를 입력해 주세요.",
    };
  }

  if (!simpleEmailPattern.test(email)) {
    return {
      code: "INVALID_EMAIL_FORMAT",
      message: "이메일 형식이 올바르지 않습니다.",
    };
  }

  return null;
}

async function authorizedFetch(path, options = {}) {
  const headers = new Headers(options.headers || {});
  if (authState && authState.accessToken) {
    headers.set("Authorization", `Bearer ${authState.accessToken}`);
  }

  const response = await fetch(path, {
    ...options,
    headers,
  });
  const payload = await readJsonIfPossible(response);

  if (response.status === 401) {
    clearAuthState();
    renderLoggedOut();
    setStatus("status-warning", payload.message || "인증이 만료되어 로그인 화면으로 돌아갑니다.", payload);
  }

  return { response, payload };
}

async function fetchAddressDetail(id, options = {}) {
  const { silent = false } = options;
  const { response, payload } = await authorizedFetch(`/api/v1/test-addresses/${id}`);

  if (!response.ok) {
    if (!silent && response.status !== 401) {
      setStatus("status-error", payload.message || "상세 정보를 불러오지 못했습니다.", payload);
    }
    return null;
  }

  const index = authState.addresses.findIndex((address) => address.id === payload.id);
  if (index >= 0) {
    authState.addresses[index] = payload;
  } else {
    authState.addresses.push(payload);
  }
  authState.selectedAddressId = payload.id;
  renderAuthenticated();

  if (!silent) {
    setStatus("status-success", "상세 정보를 불러왔습니다.", payload);
  }

  return payload;
}

async function refreshAddresses(options = {}) {
  if (!authState) {
    return;
  }

  const { preferredId = authState.selectedAddressId, silent = false } = options;
  if (!silent) {
    setStatus("status-warning", "테스트용 메일 주소 목록을 불러오는 중입니다.", {
      request: "GET /api/v1/test-addresses",
    });
  }

  refreshAddressesButton.disabled = true;
  try {
    const { response, payload } = await authorizedFetch("/api/v1/test-addresses");
    if (!response.ok) {
      if (response.status !== 401) {
        setStatus("status-error", payload.message || "목록을 불러오지 못했습니다.", payload);
      }
      return;
    }

    authState.addresses = Array.isArray(payload.addresses) ? payload.addresses : [];
    if (authState.addresses.some((address) => address.id === preferredId)) {
      authState.selectedAddressId = preferredId;
    } else if (authState.addresses.length > 0) {
      authState.selectedAddressId = authState.addresses[0].id;
    } else {
      authState.selectedAddressId = null;
    }

    renderAuthenticated();

    if (authState.selectedAddressId) {
      await fetchAddressDetail(authState.selectedAddressId, { silent: true });
    }

    if (!silent) {
      setStatus("status-success", "주소 목록을 갱신했습니다.", payload);
    }
  } catch (error) {
    setStatus("status-error", "목록 조회 중 예기치 않은 오류가 발생했습니다.", {
      code: "NETWORK_ERROR",
      message: error instanceof Error ? error.message : "unknown error",
    });
  } finally {
    refreshAddressesButton.disabled = false;
  }
}

loginForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const loginId = loginIdInput.value.trim();
  const password = passwordInput.value;
  const validationError = validateInputs(loginId, password);

  if (validationError) {
    setStatus("status-error", validationError.message, validationError);
    return;
  }

  submitButton.disabled = true;
  setStatus("status-warning", "로그인 요청을 보내는 중입니다.", {
    loginId,
    request: "POST /api/v1/auth/login",
  });

  try {
    const response = await fetch("/api/v1/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ loginId, password }),
    });

    const payload = await response.json();

    if (response.ok) {
      authState = {
        loginId,
        accessToken: payload.accessToken,
        refreshToken: payload.refreshToken,
        accessTokenExpiresIn: payload.accessTokenExpiresIn,
        refreshTokenExpiresIn: payload.refreshTokenExpiresIn,
        tokenType: payload.tokenType,
        addresses: [],
        selectedAddressId: null,
      };
      renderAuthenticated();
      setStatus("status-success", "로그인에 성공했습니다. 주소 목록을 불러옵니다.", payload);
      await refreshAddresses({ silent: true });
      return;
    }

    if (response.status === 423) {
      setStatus("status-warning", payload.message || "계정이 잠금 상태입니다.", payload);
      return;
    }

    setStatus("status-error", payload.message || "로그인에 실패했습니다.", payload);
  } catch (error) {
    setStatus("status-error", "요청 중 예기치 않은 오류가 발생했습니다.", {
      code: "NETWORK_ERROR",
      message: error instanceof Error ? error.message : "unknown error",
    });
  } finally {
    submitButton.disabled = false;
  }
});

logoutButton.addEventListener("click", async () => {
  if (!authState) {
    renderLoggedOut();
    return;
  }

  logoutButton.disabled = true;
  setStatus("status-warning", "로그아웃 요청을 보내는 중입니다.", {
    loginId: authState.loginId,
    request: "POST /api/v1/auth/logout",
  });

  try {
    const response = await fetch("/api/v1/auth/logout", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${authState.accessToken}`,
      },
    });

    const payload = await readJsonIfPossible(response);

    if (response.status === 204) {
      clearAuthState();
      renderLoggedOut();
      setStatus("status-success", "로그아웃되었습니다.", {
        status: 204,
      });
      return;
    }

    if (response.status === 401) {
      clearAuthState();
      renderLoggedOut();
      setStatus("status-warning", payload.message || "인증이 유효하지 않아 로그인 화면으로 이동합니다.", payload);
      return;
    }

    setStatus("status-error", payload.message || "로그아웃에 실패했습니다.", payload);
  } catch (error) {
    setStatus("status-error", "로그아웃 요청 중 예기치 않은 오류가 발생했습니다.", {
      code: "NETWORK_ERROR",
      message: error instanceof Error ? error.message : "unknown error",
    });
  } finally {
    logoutButton.disabled = false;
  }
});

refreshAddressesButton.addEventListener("click", async () => {
  await refreshAddresses();
});

addressList.addEventListener("click", async (event) => {
  const button = event.target.closest("button[data-address-id]");
  if (!button) {
    return;
  }

  const addressId = Number(button.dataset.addressId);
  if (!Number.isInteger(addressId)) {
    return;
  }

  await fetchAddressDetail(addressId);
});

generateCandidateButton.addEventListener("click", async () => {
  if (!authState) {
    setStatus("status-warning", "먼저 로그인해 주세요.", {});
    return;
  }

  generateCandidateButton.disabled = true;
  setStatus("status-warning", "후보 메일 주소를 생성하는 중입니다.", {
    request: "POST /api/v1/test-addresses/generate",
  });

  try {
    const { response, payload } = await authorizedFetch("/api/v1/test-addresses/generate", {
      method: "POST",
    });
    if (!response.ok) {
      if (response.status !== 401) {
        setStatus("status-error", payload.message || "후보값을 가져오지 못했습니다.", payload);
      }
      return;
    }

    createEmailInput.value = payload.email || "";
    candidateHelper.textContent = `추천 후보값을 입력창에 채웠습니다: ${payload.email}`;
    setStatus("status-success", "후보 메일 주소를 받아왔습니다.", payload);
  } catch (error) {
    setStatus("status-error", "후보 생성 요청 중 예기치 않은 오류가 발생했습니다.", {
      code: "NETWORK_ERROR",
      message: error instanceof Error ? error.message : "unknown error",
    });
  } finally {
    generateCandidateButton.disabled = false;
  }
});

createAddressForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  if (!authState) {
    setStatus("status-warning", "먼저 로그인해 주세요.", {});
    return;
  }

  const email = createEmailInput.value.trim();
  const validationError = validateCreateEmail(email);
  if (validationError) {
    setStatus("status-error", validationError.message, validationError);
    return;
  }

  createAddressButton.disabled = true;
  setStatus("status-warning", "테스트용 메일 주소를 생성하는 중입니다.", {
    request: "POST /api/v1/test-addresses",
    email,
  });

  try {
    const { response, payload } = await authorizedFetch("/api/v1/test-addresses", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ email }),
    });
    if (!response.ok) {
      if (response.status !== 401) {
        setStatus("status-error", payload.message || "주소를 생성하지 못했습니다.", payload);
      }
      return;
    }

    createEmailInput.value = "";
    candidateHelper.textContent = "후보값을 받아 입력창을 채우거나, 원하는 주소를 직접 입력하세요.";
    if (!authState.addresses) {
      authState.addresses = [];
    }
    authState.addresses.unshift(payload);
    authState.selectedAddressId = payload.id;
    renderAuthenticated();
    await fetchAddressDetail(payload.id, { silent: true });
    setStatus("status-success", "테스트용 메일 주소를 생성했습니다.", payload);
  } catch (error) {
    setStatus("status-error", "주소 생성 요청 중 예기치 않은 오류가 발생했습니다.", {
      code: "NETWORK_ERROR",
      message: error instanceof Error ? error.message : "unknown error",
    });
  } finally {
    createAddressButton.disabled = false;
  }
});

renderLoggedOut();
