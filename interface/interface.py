import requests  # библиотека для HTTP запросов
import json  # работа с JSON данными


def get_order_info(order_id):
    """
    Функция для получения информации о заказе по его ID.

    Args:
    - order_id (str): ID заказа для запроса информации.

    Returns:
    - dict: Информация о заказе в формате словаря, если запрос успешен.
            None, если произошла ошибка при запросе.
    """
    url = f'http://localhost:8000/orders/{order_id}'
    try:
        response = requests.get(url)  # отправляем GET запрос
        if response.status_code == 200:
            order_data = response.json()  # декодируем JSON ответ
            return order_data
        else:
            print(f'Ошибка получения заказа. Статус: {response.status_code}')
            return None
    except requests.exceptions.RequestException as e:
        print(f'Ошибка получения заказа: {e}')
        return None


def print_order_info(order_data):
    """
    Функция для печати информации о заказе.

    Args:
    - order_data (dict): Информация о заказе в виде словаря.
    """
    print("Информация о заказе:")
    print(f"Order UID: {order_data['order_uid']}")
    print(f"Track Number: {order_data['track_number']}")
    print(f"Entry: {order_data['entry']}")
    print("Delivery:")
    print(f"  Name: {order_data['delivery']['name']}")
    print(f"  Phone: {order_data['delivery']['phone']}")
    print(f"  Zip: {order_data['delivery']['zip']}")
    print(f"  City: {order_data['delivery']['city']}")
    print(f"  Address: {order_data['delivery']['address']}")
    print(f"  Region: {order_data['delivery']['region']}")
    print(f"  Email: {order_data['delivery']['email']}")
    print("Payment:")
    print(f"  Transaction: {order_data['payment']['transaction']}")
    print(f"  Currency: {order_data['payment']['currency']}")
    print(f"  Provider: {order_data['payment']['provider']}")
    print(f"  Amount: {order_data['payment']['amount']}")
    print(f"  Payment Date: {order_data['payment']['payment_dt']}")
    print(f"  Bank: {order_data['payment']['bank']}")
    print(f"  Delivery Cost: {order_data['payment']['delivery_cost']}")
    print(f"  Goods Total: {order_data['payment']['goods_total']}")
    print(f"  Custom Fee: {order_data['payment']['custom_fee']}")
    print("Items:")
    for item in order_data['items']:
        print(f"  - Chrt ID: {item['chrt_id']}")
        print(f"    Track Number: {item['track_number']}")
        print(f"    Price: {item['price']}")
        print(f"    Rid: {item['rid']}")
        print(f"    Name: {item['name']}")
        print(f"    Sale: {item['sale']}")
        print(f"    Size: {item['size']}")
        print(f"    Total Price: {item['total_price']}")
        print(f"    Nm ID: {item['nm_id']}")
        print(f"    Brand: {item['brand']}")
        print(f"    Status: {item['status']}")
    print(f"Locale: {order_data['locale']}")
    print(f"Internal Signature: {order_data['internal_signature']}")
    print(f"Delivery Service: {order_data['delivery_service']}")
    print(f"Shardkey: {order_data['shardkey']}")
    print(f"SM ID: {order_data['sm_id']}")
    print(f"Date Created: {order_data['date_created']}")
    print(f"OOF Shard: {order_data['oof_shard']}")


def main():
    """
    Основная функция программы. Запрашивает у пользователя ID заказа,
    получает информацию о заказе и выводит её на экран.
    """
    order_id = input('Введите ID заказа: ')
    order_data = get_order_info(order_id)
    if order_data:
        print_order_info(order_data)
    else:
        print('Не удалось получить информацию о заказе. Пожалуйста, попробуйте еще раз.')


if __name__ == '__main__':
    main()
