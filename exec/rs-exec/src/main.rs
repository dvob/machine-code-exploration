fn main() {
    let code: Vec<u8> = vec![
        0x48, 0xc7, 0xc0, 0x3c, 0x00, 0x00, 0x00, //
        0x48, 0xc7, 0xc7, 0x21, 0x00, 0x00, 0x00, //
        0x0f, 0x05, //
    ];

    unsafe {
        let a = libc::mmap(
            std::ptr::null_mut(),
            code.len(),
            libc::PROT_EXEC | libc::PROT_READ | libc::PROT_WRITE,
            libc::MAP_PRIVATE | libc::MAP_ANONYMOUS,
            -1,
            0,
        );

        std::ptr::copy_nonoverlapping(code.as_ptr(), a as *mut u8, code.len());

        let fn_ptr: fn() = std::mem::transmute(a);
        fn_ptr();
    }
}
