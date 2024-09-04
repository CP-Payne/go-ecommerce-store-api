INSERT INTO categories (id, name, description)
VALUES
('9b682f9c-c1ac-4242-98ea-f05d5b2d680c','Laptops', 'Various models of laptops and notebooks'),
('8bc5a2ca-b3ed-4fad-91b5-aaf302e070bf','Smartphones', 'Latest smartphones from top brands'),
('87785cc5-4c62-441f-8fb5-add44b598c8c','Tech Accessories', 'Accessories for your tech devices'),
('2411d196-8e64-4161-894c-2ee25fab063f','Light Gadgets', 'Innovative lighting gadgets for home and office'),
('d1d7b533-39b5-466c-8ac4-3a3bacb4e763','Wearables', 'Smartwatches, fitness trackers, and more');


INSERT INTO products (id, name, description, price, brand, sku, stock_quantity, category_id, image_url, thumbnail_url, specifications, variants)
VALUES
-- Tech Accessories
('3e2762f7-344d-4e6c-acb8-8462c67438f8',
'Wireless Mouse', 
 'Ergonomic wireless mouse with customizable buttons and DPI settings.', 
 29.99, 
 'Logitech', 
 'MOUSE-WL-LOGI', 
 150, 
 (SELECT id FROM categories WHERE name = 'Tech Accessories'), 
 'https://i.pinimg.com/originals/f5/ff/7c/f5ff7c038681a83f68144042e995aac2.jpg', 
 'https://i.pinimg.com/originals/f5/ff/7c/f5ff7c038681a83f68144042e995aac2.jpg', 
 '{"connectivity": "2.4GHz wireless", "battery_life": "12 months", "dpi": "800-1600 DPI", "color": "Black"}', 
 '[
    {"color": "Black", "price": 29.99, "sku": "MOUSE-WL-LOGI-BLK"},
    {"color": "White", "price": 29.99, "sku": "MOUSE-WL-LOGI-WHT"}
  ]'),

('749059f1-61df-4d24-934f-6179035b2149',
    'USB-C Hub', 
 '7-in-1 USB-C hub with HDMI, USB 3.0, and SD card slots.', 
 49.99, 
 'Anker', 
 'USBHUB-7IN1-ANK', 
 100, 
 (SELECT id FROM categories WHERE name = 'Tech Accessories'), 
 'https://i.pinimg.com/736x/5c/a7/c6/5ca7c606107ca60828a7c10db3d10581.jpg', 
 'https://i.pinimg.com/736x/5c/a7/c6/5ca7c606107ca60828a7c10db3d10581.jpg', 
 '{"ports": "1 HDMI, 3 USB 3.0, 1 SD card slot, 1 microSD card slot, 1 USB-C power delivery", "material": "Aluminum", "weight": "50g"}', 
 '[
    {"color": "Space Gray", "price": 49.99, "sku": "USBHUB-7IN1-ANK-GRY"},
    {"color": "Silver", "price": 49.99, "sku": "USBHUB-7IN1-ANK-SLV"}
  ]'),

-- Light Gadgets
('5b7922cd-0b27-4541-90ee-7c6230b4d90c',
    'Smart LED Bulb', 
 'Wi-Fi enabled smart LED bulb with adjustable brightness and color.', 
 19.99, 
 'Philips Hue', 
 'LED-BULB-SMART-HUE', 
 200, 
 (SELECT id FROM categories WHERE name = 'Light Gadgets'), 
 'https://i.pinimg.com/originals/1f/de/fb/1fdefbc717f4a5cc69c21d4d3d5dedb7.jpg', 
 'https://i.pinimg.com/originals/1f/de/fb/1fdefbc717f4a5cc69c21d4d3d5dedb7.jpg', 
 '{"brightness": "800 lumens", "color_temperature": "2700K-6500K", "lifespan": "25000 hours", "connectivity": "Wi-Fi, Bluetooth"}', 
 '[
    {"color": "White", "price": 19.99, "sku": "LED-BULB-SMART-HUE-WHT"},
    {"color": "Color", "price": 24.99, "sku": "LED-BULB-SMART-HUE-CLR"}
  ]'),

('b7e162b8-38a9-4bb6-b61b-68b175f1e9f5',
    'Portable LED Desk Lamp', 
 'Rechargeable LED desk lamp with touch control and adjustable arm.', 
 39.99, 
 'TaoTronics', 
 'DESK-LAMP-LED-TT', 
 75, 
 (SELECT id FROM categories WHERE name = 'Light Gadgets'), 
 'https://i.pinimg.com/originals/f7/3a/db/f73adbe044d6cd474451a0803febb3d8.jpg', 
 'https://i.pinimg.com/originals/f7/3a/db/f73adbe044d6cd474451a0803febb3d8.jpg', 
 '{"brightness_levels": "5", "color_modes": "3", "battery_life": "10 hours", "weight": "600g"}', 
 '[
    {"color": "Black", "price": 39.99, "sku": "DESK-LAMP-LED-TT-BLK"},
    {"color": "White", "price": 39.99, "sku": "DESK-LAMP-LED-TT-WHT"}
  ]'),

-- Wearables
('83b35792-bd21-4d29-bb42-0ecdf9125fb1',
    'Smartwatch X100', 
 'Feature-packed smartwatch with heart rate monitor, GPS, and more.', 
 149.99, 
 'FitGear', 
 'WATCH-SMART-X100', 
 50, 
 (SELECT id FROM categories WHERE name = 'Wearables'), 
 'https://i.pinimg.com/originals/d0/e2/a9/d0e2a9aad3d0f3f3ccc5fbbb9f2ce4fc.jpg', 
 'https://i.pinimg.com/originals/d0/e2/a9/d0e2a9aad3d0f3f3ccc5fbbb9f2ce4fc.jpg', 
 '{"display": "1.5 inch AMOLED", "battery_life": "7 days", "water_resistance": "5 ATM", "connectivity": "Bluetooth, GPS"}', 
 '[
    {"color": "Black", "band": "Silicone", "price": 149.99, "sku": "WATCH-SMART-X100-BLK-SIL"},
    {"color": "Silver", "band": "Metal", "price": 169.99, "sku": "WATCH-SMART-X100-SLV-MTL"}
  ]'),

('99734a47-a88e-48e3-9e89-c55bd515c9e2',
    'Fitness Tracker Z200', 
 'Sleek fitness tracker with sleep monitoring and step counter.', 
 69.99, 
 'HealthPlus', 
 'FITNESS-TRACK-Z200', 
 80, 
 (SELECT id FROM categories WHERE name = 'Wearables'), 
 'https://i.pinimg.com/originals/42/fc/98/42fc98825cc3372bfb94b43a69059d9b.jpg', 
 'https://i.pinimg.com/originals/42/fc/98/42fc98825cc3372bfb94b43a69059d9b.jpg', 
 '{"display": "OLED", "battery_life": "10 days", "water_resistance": "IP68", "connectivity": "Bluetooth"}', 
 '[
    {"color": "Black", "band": "Silicone", "price": 69.99, "sku": "FITNESS-TRACK-Z200-BLK"},
    {"color": "Blue", "band": "Silicone", "price": 69.99, "sku": "FITNESS-TRACK-Z200-BLU"}
  ]');

