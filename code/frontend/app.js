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

const loginIdPattern = /^[a-z0-9-]{4,32}$/;
const passwordPattern = /^(?=.*[A-Za-z])(?=.*[0-9])[^\s]{8,64}$/;
let authState = null;

function setStatus(kind, message, payload) {
  statusPanel.className = `status-panel ${kind}`;
  statusMessage.textContent = message;
  statusPayload.textContent = JSON.stringify(payload, null, 2);
}

function renderLoggedOut() {
  loginView.classList.remove("hidden");
  authView.classList.add("hidden");
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
      };
      renderAuthenticated();
      setStatus("status-success", "로그인에 성공했습니다.", payload);
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

renderLoggedOut();
