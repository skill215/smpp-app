#!/usr/bin/env python3

"""
SMPP Application Test Data Generator
==================================

This script generates test data for the SMPP application, including:
1. Text messages (text.txt) - English by default, with optional multilingual support
2. Random URLs (url.txt)

Language Options:
--------------
- Default: English only
- With --multilingual flag:
  - English
  - Chinese (简体中文)
  - Japanese (日本語)
  - Korean (한국어)
  - Hebrew (עברית)

Message Types:
------------
- Welcome messages
- Verification codes
- Special offers
- Order confirmations
- News updates
- Appointment reminders
- Security alerts
- Point rewards
- System updates
- Payment notifications

URL Generation:
-------------
- Random subdomain, domain, and TLD
- Supports http and https protocols
- Random path and query parameters
- Common TLDs: .com, .org, .net, .io, .dev

Usage Examples:
-------------
1. Generate English messages (default):
   ./generate_test_data.py

2. Generate multilingual messages:
   ./generate_test_data.py --multilingual

3. Generate custom number of messages and URLs:
   ./generate_test_data.py --msgnum 200 --urlnum 100 [--multilingual]

4. Show this help message:
   ./generate_test_data.py --help

Output Files:
-----------
- data/text.txt: Contains generated messages (English or multilingual)
- data/url.txt: Contains generated random URLs
"""

import argparse
import random
import string
import os

def generate_random_string(length, include_unicode=False):
    """Generate a random string of specified length."""
    if include_unicode:
        # Include Chinese, Japanese, Korean, and other Unicode characters
        unicode_ranges = [
            (0x4E00, 0x9FFF),   # Chinese characters
            (0x3040, 0x309F),   # Hiragana
            (0x30A0, 0x30FF),   # Katakana
            (0xAC00, 0xD7AF),   # Korean Hangul
            (0x0400, 0x04FF),   # Cyrillic
            (0x0900, 0x097F),   # Devanagari
            (0x0600, 0x06FF),   # Arabic
            (0x0590, 0x05FF),   # Hebrew
        ]
        chars = []
        for _ in range(length):
            range_start, range_end = random.choice(unicode_ranges)
            chars.append(chr(random.randint(range_start, range_end)))
        return ''.join(chars)
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

def generate_random_url():
    """Generate a random URL."""
    protocols = ['http', 'https']
    tlds = ['com', 'org', 'net', 'io', 'dev']
    
    protocol = random.choice(protocols)
    subdomain = generate_random_string(5)
    domain = generate_random_string(8)
    tld = random.choice(tlds)
    path = '/' + generate_random_string(10)
    query = '?' + generate_random_string(5) + '=' + generate_random_string(5)
    
    return f"{protocol}://{subdomain}.{domain}.{tld}{path}{query}"

def generate_random_message(multilingual=False):
    """Generate a random message."""
    if multilingual:
        templates = [
            # English templates
            "Welcome {name}! Your account has been verified.",
            "Your verification code is {code}. Valid for 5 minutes.",
            "Special offer: {discount}% off on all items! Use code {code}",
            "Your order #{order} has been confirmed.",
            "Breaking news: {news}!",
            "Reminder: Your appointment is scheduled for {time}",
            "Security alert: New login from {location}",
            "Congratulations! You've earned {points} bonus points",
            "System update: {feature} is now available",
            "Payment of ${amount} received. Transaction ID: {id}",
            # Chinese templates
            "欢迎 {name}！您的账户已验证。",
            "您的验证码是 {code}，5分钟内有效。",
            "特别优惠：全场商品{discount}%折扣！使用代码 {code}",
            "订单 #{order} 已确认。",
            "重要通知：{news}！",
            "提醒：您的预约时间是 {time}",
            "安全提醒：检测到新登录，位置 {location}",
            "恭喜！您获得了 {points} 积分",
            "系统更新：{feature} 功能已上线",
            "收到付款 ${amount}。交易编号：{id}",
            # Japanese templates
            "ようこそ {name}様！アカウントが確認されました。",
            "認証コード：{code}（有効期限5分）",
            "特別セール：全品{discount}%オフ！コード：{code}",
            "注文番号 #{order} が確認されました。",
            "お知らせ：{news}",
            "リマインダー：予約時間は {time} です",
            "セキュリティ警告：新規ログイン場所 {location}",
            "おめでとう！{points} ポイントを獲得しました",
            "システム更新：{feature} が利用可能になりました",
            "支払い完了：${amount}。取引ID：{id}",
            # Korean templates
            "환영합니다 {name}님! 계정이 인증되었습니다.",
            "인증번호는 {code}입니다. 5분 동안 유효합니다.",
            "특별 할인: 전체 상품 {discount}% 할인! 코드: {code}",
            "주문번호 #{order} 확인되었습니다.",
            "주요 소식: {news}!",
            "알림: 예약 시간은 {time}입니다",
            "보안 경고: 새로운 로그인 위치 {location}",
            "축하합니다! {points} 포인트를 획득하셨습니다",
            "시스템 업데이트: {feature} 기능이 추가되었습니다",
            "결제 완료: ${amount}. 거래 ID: {id}",
            # Hebrew templates
            "!ברוך הבא {name}! החשבון שלך אומת",
            "קוד האימות שלך הוא {code}. תקף ל-5 דקות",
            "מבצע מיוחד: {discount}% הנחה על כל הפריטים! השתמש בקוד {code}",
            "הזמנה מספר #{order} אושרה",
            "חדשות חשובות: {news}!",
            "תזכורת: הפגישה שלך נקבעה ל-{time}",
            "התראת אבטחה: כניסה חדשה מ-{location}",
            "מזל טוב! צברת {points} נקודות בונוס",
            "עדכון מערכת: {feature} זמין עכשיו",
            "התקבל תשלום בסך ${amount}. מזהה עסקה: {id}"
        ]
    else:
        templates = [
            # English templates only
            "Welcome {name}! Your account has been verified.",
            "Your verification code is {code}. Valid for 5 minutes.",
            "Special offer: {discount}% off on all items! Use code {code}",
            "Your order #{order} has been confirmed.",
            "Breaking news: {news}!",
            "Reminder: Your appointment is scheduled for {time}",
            "Security alert: New login from {location}",
            "Congratulations! You've earned {points} bonus points",
            "System update: {feature} is now available",
            "Payment of ${amount} received. Transaction ID: {id}"
        ]
    
    template = random.choice(templates)
    
    # Replace placeholders with random values
    replacements = {
        '{name}': generate_random_string(8, include_unicode=multilingual),
        '{code}': ''.join(random.choices(string.digits, k=6)),
        '{discount}': str(random.randint(10, 70)),
        '{order}': ''.join(random.choices(string.digits, k=8)),
        '{news}': ' '.join(generate_random_string(6, include_unicode=multilingual) for _ in range(2)),
        '{time}': f"{random.randint(1, 12)}:{random.randint(0, 59):02d} {random.choice(['AM', 'PM'])}",
        '{location}': f"{generate_random_string(8, include_unicode=multilingual)}",
        '{points}': str(random.randint(100, 1000)),
        '{feature}': generate_random_string(10, include_unicode=multilingual),
        '{amount}': f"{random.randint(10, 1000)}.{random.randint(0, 99):02d}",
        '{id}': generate_random_string(12).upper()
    }
    
    for key, value in replacements.items():
        if key in template:
            template = template.replace(key, value)
    
    return template

def main():
    parser = argparse.ArgumentParser(
        description='Generate test data for SMPP application',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    parser.add_argument(
        '--msgnum',
        type=int,
        default=100,
        help='Number of messages to generate (default: 100)'
    )
    parser.add_argument(
        '--urlnum',
        type=int,
        default=50,
        help='Number of URLs to generate (default: 50)'
    )
    parser.add_argument(
        '--multilingual',
        action='store_true',
        help='Enable multilingual message generation (default: English only)'
    )
    args = parser.parse_args()
    
    # Create data directory if it doesn't exist
    os.makedirs('data', exist_ok=True)
    
    # Generate messages
    lang_mode = "multiple languages" if args.multilingual else "English"
    print(f"Generating {args.msgnum} messages in {lang_mode}...")
    with open('data/text.txt', 'w', encoding='utf-8') as f:
        for _ in range(args.msgnum):
            f.write(generate_random_message(args.multilingual) + '\n')
    
    # Generate URLs
    print(f"Generating {args.urlnum} random URLs...")
    with open('data/url.txt', 'w', encoding='utf-8') as f:
        for _ in range(args.urlnum):
            f.write(generate_random_url() + '\n')
    
    print("\nDone! Files generated:")
    print("- data/text.txt (Messages in " + lang_mode + ")")
    print("- data/url.txt (Random URLs)")
    
    if args.multilingual:
        print("\nSupported languages:")
        print("- English")
        print("- Chinese (简体中文)")
        print("- Japanese (日本語)")
        print("- Korean (한국어)")
        print("- Hebrew (עברית)")

if __name__ == '__main__':
    main() 