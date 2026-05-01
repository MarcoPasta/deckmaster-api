<script setup lang="ts">
import { ref } from "vue";

const firstName = ref("");
const lastName = ref("");
const currentPassword = ref("");
const newPassword = ref("");
const confirmPassword = ref("");
const profileImage = ref<string | null>(null);
const showCurrentPw = ref(false);
const showNewPw = ref(false);
const showConfirmPw = ref(false);

const fileInput = ref<HTMLInputElement | null>(null);

function onFileSelect(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = (e) => {
    profileImage.value = e.target?.result as string;
  };
  reader.readAsDataURL(file);
}

function removeImage() {
  profileImage.value = null;
  if (fileInput.value) fileInput.value.value = "";
}

function saveProfile() {
  // TODO: call API to persist profile changes
  console.log("Saving profile:", {
    firstName: firstName.value,
    lastName: lastName.value,
  });
}

function savePassword() {
  if (newPassword.value !== confirmPassword.value) return;
  // TODO: call API to change password
  console.log("Changing password");
  currentPassword.value = "";
  newPassword.value = "";
  confirmPassword.value = "";
}
</script>

<template>
  <div class="profile-page">
    <div class="profile-card">
      <h2 class="section-title">Profile</h2>

      <!-- Avatar -->
      <div class="avatar-section">
        <div class="avatar" @click="fileInput?.click()">
          <img
            v-if="profileImage"
            :src="profileImage"
            alt="Profile"
            class="avatar-img"
          />
          <span v-else class="avatar-placeholder">+</span>
        </div>
        <div class="avatar-actions">
          <v-btn size="small" variant="outlined" @click="fileInput?.click()">
            Upload Photo
          </v-btn>
          <v-btn
            v-if="profileImage"
            size="small"
            variant="text"
            color="error"
            @click="removeImage"
          >
            Remove
          </v-btn>
        </div>
        <input
          ref="fileInput"
          type="file"
          accept="image/*"
          hidden
          @change="onFileSelect"
        />
      </div>

      <v-divider class="my-4" />

      <!-- Name fields -->
      <div class="form-row">
        <v-text-field
          v-model="firstName"
          label="First Name"
          variant="outlined"
          density="compact"
          hide-details
        />
        <v-text-field
          v-model="lastName"
          label="Last Name"
          variant="outlined"
          density="compact"
          hide-details
        />
      </div>

      <div class="form-actions">
        <v-btn color="deep-purple-accent-2" @click="saveProfile">
          Save Changes
        </v-btn>
      </div>

      <v-divider class="my-6" />

      <!-- Password section -->
      <h2 class="section-title">Change Password</h2>

      <div class="form-stack">
        <v-text-field
          v-model="currentPassword"
          label="Current Password"
          :type="showCurrentPw ? 'text' : 'password'"
          variant="outlined"
          density="compact"
          hide-details
          :append-inner-icon="showCurrentPw ? 'mdi-eye-off' : 'mdi-eye'"
          @click:append-inner="showCurrentPw = !showCurrentPw"
        />
        <v-text-field
          v-model="newPassword"
          label="New Password"
          :type="showNewPw ? 'text' : 'password'"
          variant="outlined"
          density="compact"
          hide-details
          :append-inner-icon="showNewPw ? 'mdi-eye-off' : 'mdi-eye'"
          @click:append-inner="showNewPw = !showNewPw"
        />
        <v-text-field
          v-model="confirmPassword"
          label="Confirm New Password"
          :type="showConfirmPw ? 'text' : 'password'"
          variant="outlined"
          density="compact"
          hide-details
          :error="confirmPassword.length > 0 && newPassword !== confirmPassword"
          :append-inner-icon="showConfirmPw ? 'mdi-eye-off' : 'mdi-eye'"
          @click:append-inner="showConfirmPw = !showConfirmPw"
        />
      </div>

      <div class="form-actions">
        <v-btn
          color="deep-purple-accent-2"
          :disabled="
            !currentPassword || !newPassword || newPassword !== confirmPassword
          "
          @click="savePassword"
        >
          Update Password
        </v-btn>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-page {
  display: flex;
  justify-content: center;
  padding: 24px 16px;
  box-sizing: border-box;
}

.profile-card {
  width: 100%;
  max-width: 560px;
  background: var(--code-bg);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 28px 32px;
  box-shadow: var(--shadow);
}

.section-title {
  font-size: 20px;
  margin: 0 0 16px;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 2px dashed var(--border);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  cursor: pointer;
  flex-shrink: 0;
  transition: border-color 0.2s;
}

.avatar:hover {
  border-color: var(--accent-border);
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  font-size: 28px;
  color: var(--text);
  opacity: 0.5;
}

.avatar-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.form-row {
  display: flex;
  gap: 12px;
  margin-bottom: 4px;
}

.form-stack {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.my-4 {
  margin: 16px 0;
}

.my-6 {
  margin: 24px 0;
}
</style>
