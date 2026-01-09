mod git;

use tauri::{
    menu::{Menu, MenuItem},
    tray::{TrayIconBuilder, TrayIconEvent},
    Manager,
};

pub fn run() {
    tauri::Builder::default()
        .setup(|app| {
            let quit_i = MenuItem::with_id(app, "quit", "Quit", true, None::<&str>)?;
            let show_i = MenuItem::with_id(app, "show", "Show", true, None::<&str>)?;
            let hide_i = MenuItem::with_id(app, "hide", "Hide", true, None::<&str>)?;
            let menu = Menu::with_items(app, &[&show_i, &hide_i, &quit_i])?;
            let _tray = TrayIconBuilder::new()
                .icon(app.default_window_icon().unwrap().clone())
                .menu(&menu)
                .show_menu_on_left_click(false)
                .on_menu_event(|app, event| match event.id.as_ref() {
                    "quit" => { app.exit(0); }
                    "show" => { if let Some(window) = app.get_webview_window("main") { let _ = window.show(); let _ = window.set_focus(); } }
                    "hide" => { if let Some(window) = app.get_webview_window("main") { let _ = window.hide(); } }
                    _ => { println!("menu item {:?} not handled", event.id); }
                })
                .on_tray_icon_event(|tray, event| {
                    if let TrayIconEvent::Click { button: tauri::tray::MouseButton::Left, .. } = event {
                        let app = tray.app_handle();
                        if let Some(window) = app.get_webview_window("main") { let _ = window.show(); let _ = window.set_focus(); }
                    }
                })
                .build(app)?;
            Ok(())
        })
        .plugin(tauri_plugin_fs::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![
            git::get_git_branches,
            git::get_git_commits,
            git::get_commit_files,
            git::get_commit_file_diff,
            git::get_git_remotes,
            git::get_git_remote_branches,
            git::get_git_tags,
            git::get_git_stashes,
            git::git_checkout_remote_branch,
            git::switch_branch,
            git::get_ssh_key_info,
            git::generate_ssh_key,
            git::get_git_status,
            git::get_git_diff,
            git::get_avatar,
            git::clear_avatar_cache,
            git::git_commit,
            git::git_stage_file,
            git::git_stage_all,
            git::git_unstage_file,
            git::git_unstage_all,
            git::git_fetch,
            git::git_pull,
            git::git_push,
            git::git_create_branch,
            git::git_create_tag,
            git::git_delete_branch,
            git::git_discard_changes,
            git::git_stash_save,
            git::git_stash_drop,
            git::git_stash_pop
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
