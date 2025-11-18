// 认证页面功能
class AuthPage {
  constructor() {
    // 从路径中获取tab值
    const path = window.location.pathname;
    const pathParts = path.split('/');
    const lastPath = pathParts[pathParts.length - 1];

    // 匹配login或register，否则默认为login
    this.currentTab = (lastPath === 'login' || lastPath === 'register') ? lastPath : 'login';
    this.init();
  }

  init() {
    this.initTabs();
    this.initPasswordToggle();
    this.initFormValidation();
    // 根据URL参数设置初始标签页
    this.switchTab(this.currentTab);
  }

  // 初始化标签页切换
  initTabs() {
    const tabBtns = document.querySelectorAll('.tab-btn');
    const authForms = document.querySelectorAll('.auth-form');

    tabBtns.forEach(btn => {
      btn.addEventListener('click', () => {
        const targetTab = btn.dataset.tab;
        this.switchTab(targetTab);
      });
    });
  }

  switchTab(tabName) {
    // 更新按钮状态
    document.querySelectorAll('.tab-btn').forEach(btn => {
      btn.classList.toggle('active', btn.dataset.tab === tabName);
    });

    // 更新表单显示
    document.querySelectorAll('.auth-form').forEach(form => {
      form.classList.toggle('active', form.id === `${tabName}Form`);
    });

    this.currentTab = tabName;

    // 更新URL参数
    // 更新URL路径
    const basePath = window.location.pathname.substring(0, window.location.pathname.lastIndexOf('/'));
    const newPath = basePath + '/' + tabName;
    window.history.pushState({}, '', newPath);

    // 重置表单
    this.resetFormValidation();
  }

  // 初始化密码显示/隐藏
  initPasswordToggle() {
    document.querySelectorAll('.toggle-password').forEach(btn => {
      btn.addEventListener('click', (e) => {
        const targetId = e.target.dataset.target;
        const input = document.getElementById(targetId);
        const type = input.type === 'password' ? 'text' : 'password';
        input.type = type;
        e.target.textContent = type === 'password' ? '👁️' : '🔒';
      });
    });
  }

  // 初始化表单验证
  initFormValidation() {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');

    loginForm.addEventListener('submit', (e) => this.handleLogin(e));
    registerForm.addEventListener('submit', (e) => this.handleRegister(e));

    // 实时验证
    this.initRealTimeValidation();
  }

  initRealTimeValidation() {
    // 登录表单实时验证
    const loginEmail = document.getElementById('loginEmail');
    const loginPassword = document.getElementById('loginPassword');

    loginEmail.addEventListener('input', () => this.validateEmail(loginEmail, 'loginEmailError'));
    loginPassword.addEventListener('input', () => this.validatePassword(loginPassword, 'loginPasswordError'));

    // 注册表单实时验证
    const registerUsername = document.getElementById('registerUsername');
    const registerEmail = document.getElementById('registerEmail');
    const registerPassword = document.getElementById('registerPassword');
    const confirmPassword = document.getElementById('confirmPassword');

    registerUsername.addEventListener('input', () => this.validateUsername(registerUsername, 'registerUsernameError'));
    registerEmail.addEventListener('input', () => this.validateEmail(registerEmail, 'registerEmailError'));
    registerPassword.addEventListener('input', () => {
      this.validatePassword(registerPassword, 'registerPasswordError');
      this.updatePasswordStrength(registerPassword.value);
    });
    confirmPassword.addEventListener('input', () => this.validateConfirmPassword(registerPassword, confirmPassword, 'confirmPasswordError'));
  }

  // 表单验证方法
  validateEmail(input, errorId) {
    const email = input.value.trim();
    const errorElement = document.getElementById(errorId);

    if (!email) {
      this.showError(input, errorElement, '邮箱不能为空');
      return false;
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
      this.showError(input, errorElement, '请输入有效的邮箱地址');
      return false;
    }

    this.clearError(input, errorElement);
    return true;
  }

  validatePassword(input, errorId) {
    const password = input.value;
    const errorElement = document.getElementById(errorId);

    if (!password) {
      this.showError(input, errorElement, '密码不能为空');
      return false;
    }

    if (password.length < 6) {
      this.showError(input, errorElement, '密码至少需要6位字符');
      return false;
    }

    this.clearError(input, errorElement);
    return true;
  }

  validateUsername(input, errorId) {
    const username = input.value.trim();
    const errorElement = document.getElementById(errorId);

    if (!username) {
      this.showError(input, errorElement, '用户名不能为空');
      return false;
    }

    if (username.length < 3) {
      this.showError(input, errorElement, '用户名至少需要3位字符');
      return false;
    }

    if (!/^[a-zA-Z0-9_]+$/.test(username)) {
      this.showError(input, errorElement, '用户名只能包含字母、数字和下划线');
      return false;
    }

    this.clearError(input, errorElement);
    return true;
  }

  validateConfirmPassword(passwordInput, confirmInput, errorId) {
    const password = passwordInput.value;
    const confirmPassword = confirmInput.value;
    const errorElement = document.getElementById(errorId);

    if (!confirmPassword) {
      this.showError(confirmInput, errorElement, '请确认密码');
      return false;
    }

    if (password !== confirmPassword) {
      this.showError(confirmInput, errorElement, '两次输入的密码不一致');
      return false;
    }

    this.clearError(confirmInput, errorElement);
    return true;
  }

  showError(input, errorElement, message) {
    input.classList.add('error');
    errorElement.textContent = message;
  }

  clearError(input, errorElement) {
    input.classList.remove('error');
    errorElement.textContent = '';
  }

  // 密码强度检测
  updatePasswordStrength(password) {
    const strengthBar = document.getElementById('passwordStrength');
    const strengthText = document.getElementById('passwordStrengthText');

    let strength = 0;
    let text = '密码强度';

    if (password.length >= 6) strength += 1;
    if (password.length >= 8) strength += 1;
    if (/[A-Z]/.test(password)) strength += 1;
    if (/[0-9]/.test(password)) strength += 1;
    if (/[^A-Za-z0-9]/.test(password)) strength += 1;

    strengthBar.className = 'strength-fill';

    if (password.length === 0) {
      text = '密码强度';
    } else if (strength <= 2) {
      strengthBar.classList.add('weak');
      text = '弱';
    } else if (strength <= 4) {
      strengthBar.classList.add('medium');
      text = '中';
    } else {
      strengthBar.classList.add('strong');
      text = '强';
    }

    strengthText.textContent = text;
  }

  // 表单提交处理
  handleLogin(e) {
    e.preventDefault();

    const email = document.getElementById('loginEmail');
    const password = document.getElementById('loginPassword');
    const remember = document.querySelector('input[name="remember"]');
    const isEmailValid = this.validateEmail(email, 'loginEmailError');
    const isPasswordValid = this.validatePassword(password, 'loginPasswordError');

    if (isEmailValid && isPasswordValid) {
      // 提交数据到/login路由
      fetch('/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
          'email': email.value,
          'password': password.value,
          'remember': remember.checked ? 'on' : ''
        })
      })
          .then(response => response.json())
          .then(data => {
            console.log('登录响应:', data);
            if (data.status === 'success') {
              this.showSuccess('登录成功！正在跳转...');
              setTimeout(() => {
                window.location.href = '/';
              }, 500);
            } else {
              this.showError(email, document.getElementById('loginEmailError'), data.message || '登录失败');
            }
          })
          .catch(error => {
            console.error('登录错误:', error);
            this.showError(email, document.getElementById('loginEmailError'), '网络错误，请稍后重试');
          });
    }
  }

  handleRegister(e) {
    e.preventDefault();

    const username = document.getElementById('registerUsername');
    const email = document.getElementById('registerEmail');
    const password = document.getElementById('registerPassword');
    const confirmPassword = document.getElementById('confirmPassword');
    const agreeTerms = document.querySelector('input[name="agreeTerms"]');

    const isUsernameValid = this.validateUsername(username, 'registerUsernameError');
    const isEmailValid = this.validateEmail(email, 'registerEmailError');
    const isPasswordValid = this.validatePassword(password, 'registerPasswordError');
    const isConfirmValid = this.validateConfirmPassword(password, confirmPassword, 'confirmPasswordError');

    if (!agreeTerms.checked) {
      this.showError(agreeTerms, document.createElement('div'), '请同意服务条款和隐私政策');
      return;
    }

    if (isUsernameValid && isEmailValid && isPasswordValid && isConfirmValid) {
      // 模拟注册成功
      // 提交数据到/login路由
      fetch('/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({
          'username': username.value,
          'email': email.value,
          'password': password.value,
          'confirmPassword': confirmPassword.value,
          'agreeTerms': agreeTerms.checked ? 'on' : ''
        })
      })
          .then(response => response.json())
          .then(data => {
            console.log('注册响应:', data);
            if (data.status === 'success') {
              this.showSuccess('注册成功！正在跳转...');
              setTimeout(() => {
                this.switchTab('login');
              }, 1500);
            } else {
              // 优化错误消息显示
              let errorMessage = '注册失败';
              if (data.message) {
                // 根据不同的错误类型显示更友好的提示
                if (data.message.includes('邮箱')) {
                  errorMessage = data.message;
                  this.showError(email, document.getElementById('registerEmailError'), errorMessage);
                } else if (data.message.includes('用户名')) {
                  errorMessage = data.message;
                  this.showError(username, document.getElementById('registerUsernameError'), errorMessage);
                } else if (data.message.includes('密码')) {
                  errorMessage = data.message;
                  this.showError(password, document.getElementById('registerPasswordError'), errorMessage);
                } else {
                  // 默认显示在邮箱错误区域
                  errorMessage = data.message;
                  this.showError(email, document.getElementById('registerEmailError'), errorMessage);
                }
              } else {
                this.showError(email, document.getElementById('registerEmailError'), '注册失败，请稍后重试');
              }
            }
          })
          .catch(error => {
            console.error('注册错误:', error);
            this.showError(email, document.getElementById('registerUsernameError'), '网络错误，请稍后重试');
          });
    }
  }

  showSuccess(message) {
    // 在实际应用中可以使用更美观的提示组件
    // alert(message);
    customAlert.success(message, 1500);
  }

  resetFormValidation() {
    document.querySelectorAll('.error-message').forEach(el => {
      el.textContent = '';
    });
    document.querySelectorAll('input.error').forEach(el => {
      el.classList.remove('error');
    });
  }
}

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', () => {
  new AuthPage();
});
