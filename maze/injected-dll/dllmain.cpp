// dllmain.cpp : Defines the entry point for the DLL application.
#include "pch.h"
#include <iostream>
#include <string>
#include "mem.h"
#include "cppdump/Assembly_CSharp/Assembly_CSharp.hpp"

std::string str(DLL2SDK::mscorlib::System::String* in) {
	std::string out = "";
	for (int i = 0; i < in->m_stringLength_; i++) {
		char c = *(char*)(&in->m_firstChar_ + i);
		out.push_back(c);
	}
	return out;
}

std::string xyz(char* label, float x, float y, float z) {
	char  buffer[100];
	// sprintf_s(buffer, 100, "%20s: %7.3f %7.3f %7.3f", label, x, y, z);
	sprintf_s(buffer, 100, " { %7.3f, %7.3f, %7.3f },", x, y, z);
	return buffer;
}

std::string xyz(char* label, DLL2SDK::UnityEngine_CoreModule::UnityEngine::Vector3 vec3) {
	return xyz(label, vec3.x_, vec3.y_, vec3.z_);
}

void debug(DLL2SDK::Assembly_CSharp::ServerManager* serverManager) {
	std::cout << xyz("Position", serverManager->position_.x_, serverManager->position_.y_, serverManager->position_.z_) << std::endl;
	std::cout << xyz("Current position", serverManager->current_position_.x_, serverManager->current_position_.y_, serverManager->current_position_.z_) << std::endl;
	std::cout << xyz("Teleport player", serverManager->teleport_player_x_, serverManager->teleport_player_y_, serverManager->teleport_player_z_) << std::endl;
	std::cout << xyz("Teleport", serverManager->teleport_x_, serverManager->teleport_y_, serverManager->teleport_z_) << std::endl;
	std::cout << "================================================================" << std::endl;
}

DLL2SDK::UnityEngine_CoreModule::UnityEngine::Vector3 position(DLL2SDK::UnityEngine_CoreModule::UnityEngine::GameObject* obj) {
	return obj->get_transform_1341()->get_position_1631();
}

void teleport(DLL2SDK::Assembly_CSharp::ServerManager* serverManager, float x, float y, float z) {
	serverManager->teleport_x_ = x;
	serverManager->teleport_y_ = y;
	serverManager->teleport_z_ = z;
	serverManager->teleportForward_359();
}

void printHighscores(DLL2SDK::Assembly_CSharp::ServerManager* serverManager) {

	for (int i = 0; i < serverManager->racemanager_->highscore_->t_personal_highscores_->Length; i++) {
		DLL2SDK::Unity_TextMeshPro::TMPro::TMP_Text* hs = serverManager->racemanager_->highscore_->t_personal_highscores_->Items[i];
		std::cout << hs->m_text_ << std::endl;
}

DLL2SDK::Assembly_CSharp::ServerManager* getServerManager(uintptr_t moduleBase) {
	uintptr_t serverManagerAddr = mem::FindDMAAddy(moduleBase + 0xAFFA48, { 0x268, 0xD20, 0x10, 0x28, 0x00 });
	return (DLL2SDK::Assembly_CSharp::ServerManager*) serverManagerAddr;
}

DWORD WINAPI RaceHack(HMODULE hModule) {
	AllocConsole();
	FILE* f;
	freopen_s(&f, "CONIN$", "r", stdin);
	freopen_s(&f, "CONOUT$", "w", stdout);

	uintptr_t moduleBase = (uintptr_t)GetModuleHandle(L"GameAssembly.dll");
	DLL2SDK::GameAssemblyBase = moduleBase;
	DLL2SDK::Assembly_CSharp::ServerManager* serverManager = nullptr;
	while (true) {
		if (GetAsyncKeyState(VK_END) & 1) {
			break;
		}
		if (GetAsyncKeyState(VK_HOME) & 1) {
            // enable syncing last_checkpoint_ with lastHearbeat_
			serverManager = getServerManager(moduleBase);
		}
		if (serverManager != nullptr) {
			serverManager->racemanager_->last_checkpoint_ = serverManager->lastHeartbeat_;
			std::cout << "Keeping race checkpoint time in sync with heartbeat: " << serverManager->racemanager_->last_checkpoint_ << " " << serverManager->lastHeartbeat_ << std::endl;
		}
		Sleep(10);
	}
	fclose(f);
	FreeConsole();
	FreeLibraryAndExitThread(hModule, 0);
	return 0;
}


DWORD WINAPI InteractiveConsole(HMODULE hModule) {
	AllocConsole();
	FILE* f;
	freopen_s(&f, "CONIN$", "r", stdin);
	freopen_s(&f, "CONOUT$", "w", stdout);


	uintptr_t moduleBase = (uintptr_t)GetModuleHandle(L"GameAssembly.dll");

    // thank you GuidedHacking, this is line is left from the tutorial
	std::cout << "Connected! OG for a fee, stay sippin' fam\n";

	std::cout << "> ";
	for (std::string line; std::getline(std::cin, line);) {
		DLL2SDK::Assembly_CSharp::ServerManager* serverManager = getServerManager(moduleBase);

		if (line == "" || line[0] == 'q') {
			std::cout << "Exiting..." << std::endl;
			break;
		}

        // changing current_position and/or position here doesn't have any effect

        // print some debugging information and toggle emojibar
		if (line[0] == 'd') {
			std::cout << "================================================================" << std::endl;
			std::cout << "Teleporters: " << serverManager->teleporters_ << std::endl;
			std::cout << "UID: " << serverManager->uid_ << std::endl;
			std::cout << "Host: " << str(serverManager->host_) << std::endl;
			std::cout << "Secret: " << str(serverManager->getSecret_364()) << std::endl;
			std::cout << "Emoji: " << serverManager->emoji_ << std::endl;
			std::cout << "================================================================" << std::endl;
            serverManager->emojibar_active_ = (serverManager->emojibar_active_ + 1) % 2;
		}

        // teleport by x y z offset
		if (line[0] == 't') {
			sscanf_s(line.c_str(), "t %f %f %f", &serverManager->teleport_x_, &serverManager->teleport_y_, &serverManager->teleport_z_);
			debug(serverManager);
			serverManager->teleportForward_359();
			debug(serverManager);
		}
		std::cout << "> ";
	}
	fclose(f);
	FreeConsole();
	FreeLibraryAndExitThread(hModule, 0);
	return 0;
}

BOOL APIENTRY DllMain(HMODULE hModule,
	DWORD  ul_reason_for_call,
	LPVOID lpReserved
)
{
	switch (ul_reason_for_call)
	{
	case DLL_PROCESS_ATTACH:
		CloseHandle(CreateThread(nullptr, 0, (LPTHREAD_START_ROUTINE)RaceHack, hModule, 0, nullptr));
	case DLL_THREAD_ATTACH:
	case DLL_THREAD_DETACH:
	case DLL_PROCESS_DETACH:
		break;
	}
	return TRUE;
}
