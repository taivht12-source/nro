package combine;

import consts.ConstNpc;
import item.Item;
import item.Item.ItemOption;
import player.Player;
import player.Service.InventoryService;
import services.ItemService;
import services.Service;
import utils.Util;

public class EpSaoTrangBi {

    // Hiển thị thông tin kết hợp trước khi nâng cấp sao trang bị
    public static void showInfoCombine(Player player) {
        if (player.combineNew.itemsCombine.size() == 2) {
            Item trangBi = null;
            Item daPhaLe = null;
            for (Item item : player.combineNew.itemsCombine) {
                if (CombineSystem.isTrangBiPhaLeHoa(item)) {
                    trangBi = item;
                } else if (CombineSystem.isDaPhaLe(item)) {
                    daPhaLe = item;
                }
            }
            int star = 0;
            int starEmpty = 0;
            if (trangBi != null && daPhaLe != null) {
                for (ItemOption io : trangBi.itemOptions) {
                    if (io.optionTemplate.id == 102) {
                        star = io.param;
                    } else if (io.optionTemplate.id == 107) {
                        starEmpty = io.param;
                    }
                }
                if (starEmpty <= 9) {
                    if (starEmpty >= 8 && !CombineService.gI().CheckSlot(trangBi, starEmpty)) {
                        CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                                "Cần cường hóa lỗ sao pha lê thứ " + (starEmpty == 8 ? "8" : "9") + " trước khi ép vào", "Đóng");
                        return;
                    }

                    player.combineNew.gemCombine = CombineSystem.getGemEpSao(star);
                    String npcSay = trangBi.template.name + "\n|2|";
                    for (ItemOption io : trangBi.itemOptions) {
                        if (io.optionTemplate.id != 102) {
                            npcSay += io.getOptionString() + "\n";
                        }
                    }

                    if (daPhaLe.template.type == 30) {
                        for (ItemOption io : daPhaLe.itemOptions) {
                            npcSay += "|7|" + io.getOptionString() + "\n";
                        }
                    } else {
                        npcSay += "|7|" + ItemService.gI().getItemOptionTemplate(CombineSystem.getOptionDaPhaLe(daPhaLe)).name
                                .replaceAll("#", CombineSystem.getParamDaPhaLe(daPhaLe) + "") + "\n";
                    }
                    npcSay += "|1|Cần " + Util.numberToMoney(player.combineNew.gemCombine) + " ngọc";
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.MENU_START_COMBINE, npcSay,
                            "Nâng cấp\ncần " + player.combineNew.gemCombine + " ngọc");
                } else {
                    CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                            "Cần 1 trang bị có lỗ sao pha lê và 1 loại đá pha lê để ép vào, và lỗ sao tối đa là 9", "Đóng");
                }
            } else {
                CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                        "Cần 1 trang bị có lỗ sao pha lê và 1 loại đá pha lê để ép vào", "Đóng");
            }
        } else {
            CombineService.gI().baHatMit.createOtherMenu(player, ConstNpc.IGNORE_MENU,
                    "Cần 1 trang bị có lỗ sao pha lê và 1 loại đá pha lê để ép vào", "Đóng");
        }
    }

    // Thực hiện việc nâng cấp sao trang bị
    public static void epSaoTrangBi(Player player) {
        // Kiểm tra xem người chơi có 2 vật phẩm cần kết hợp không
        if (player.combineNew.itemsCombine.size() == 2) {

            // Lấy số lượng ngọc cần thiết để thực hiện ép sao từ đối tượng combineNew
            int gem = player.combineNew.gemCombine;

            // Kiểm tra nếu người chơi không đủ ngọc, thì thông báo và kết thúc hàm
            if (player.inventory.gem < gem) {
                Service.gI().sendThongBao(player, "Không đủ ngọc để thực hiện");
                return;
            }

            // Khai báo hai biến vật phẩm: trang bị và đá pha lê
            Item trangBi = null;
            Item daPhaLe = null;

            // Duyệt qua các vật phẩm trong danh sách combineNew để phân biệt trang bị và đá pha lê
            for (Item item : player.combineNew.itemsCombine) {
                if (CombineSystem.isTrangBiPhaLeHoa(item)) {
                    trangBi = item;  // Nếu là trang bị, lưu vào biến trangBi
                } else if (CombineSystem.isDaPhaLe(item)) {
                    daPhaLe = item;  // Nếu là đá pha lê, lưu vào biến daPhaLe
                }
            }

            // Khởi tạo các biến sao của trang bị
            int star = 0;
            int starEmpty = 0;

            // Kiểm tra nếu có cả trang bị và đá pha lê
            if (trangBi != null && daPhaLe != null) {
                ItemOption optionStar = null;

                // Duyệt qua các tùy chọn của trang bị để lấy sao và lỗ sao pha lê
                for (ItemOption io : trangBi.itemOptions) {
                    if (io.optionTemplate.id == 102) {
                        star = io.param;  // Lấy số sao hiện tại của trang bị
                        optionStar = io;  // Lưu lại tùy chọn sao
                    } else if (io.optionTemplate.id == 107) {
                        starEmpty = io.param;  // Lấy lỗ sao pha lê
                    }
                }

                // Kiểm tra xem số sao của trang bị có nhỏ hơn lỗ sao pha lê hay không
                if (star < starEmpty) {
                    // Nếu sao nhỏ hơn lỗ sao pha lê và lỗ sao pha lê lớn hơn hoặc bằng 8,
                    // kiểm tra xem có cường hóa đúng mức hay không. Nếu không, thông báo lỗi.
                    if (starEmpty >= 8 && !CombineService.gI().CheckSlot(trangBi, starEmpty)) {
                        Service.gI().sendThongBao(player, "Cần cường hóa lỗ sao pha lê thứ " + (starEmpty == 8 ? "8" : "9") + " trước khi ép vào");
                        return;
                    }

                    // Trừ ngọc khỏi kho của người chơi
                    player.inventory.subGem(gem);

                    // Lấy các thông tin tùy chọn của đá pha lê
                    int optionId = CombineSystem.getOptionDaPhaLe(daPhaLe);
                    int param = CombineSystem.getParamDaPhaLe(daPhaLe);

                    // Kiểm tra nếu có tùy chọn tương ứng trong trang bị
                    ItemOption option = null;
                    for (ItemOption io : trangBi.itemOptions) {
                        if (io.optionTemplate.id == optionId) {
                            option = io;
                            break;
                        }
                    }

                    // Nếu sao đã cường hóa và lỗ sao pha lê là 8 hoặc 9, tạo một tùy chọn mới cho trang bị
                    if (optionStar != null && starEmpty >= 8) {
                        ItemOption newOption = new ItemOption(optionId, param);
                        trangBi.itemOptions.add(newOption);  // Thêm tùy chọn mới vào trang bị

                        // Cập nhật sao và thông báo thành công
                        if (starEmpty == 8) {
                            optionStar.param = 8;
                            Service.gI().sendThongBao(player, "Đã ép sao lên 8 thành công!");
                        } else if (starEmpty == 9) {
                            optionStar.param = 9;
                            Service.gI().sendThongBao(player, "Đã ép sao lên 9 thành công!");
                        }
                    } else {
                        // Nếu không có tùy chọn sao, hoặc sao chưa cường hóa, tăng giá trị sao và thêm tùy chọn mới
                        if (option != null) {
                            option.param += param;  // Cập nhật giá trị của tùy chọn
                        } else {
                            trangBi.itemOptions.add(new ItemOption(optionId, param));  // Thêm tùy chọn mới vào trang bị
                        }

                        // Nếu có sao, tăng sao lên
                        if (optionStar != null) {
                            optionStar.param++;
                        } else {
                            trangBi.itemOptions.add(new ItemOption(102, 1));  // Nếu không có sao, thêm tùy chọn sao mới
                        }
                    }

                    // Giảm số lượng đá pha lê trong kho
                    InventoryService.gI().subQuantityItemsBag(player, daPhaLe, 1);

                    // Gửi hiệu ứng thành công khi ép sao
                    CombineService.gI().sendEffectSuccessCombine(player);

                    // Cập nhật lại kho và thông tin tiền của người chơi
                    InventoryService.gI().sendItemBags(player);
                    Service.gI().sendMoney(player);

                    // Mở lại màn hình kết hợp
                    CombineService.gI().reOpenItemCombine(player);
                }
            }
        }
    }
}
